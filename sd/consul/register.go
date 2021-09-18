package consul

import (
	"fmt"
	"log"
	"time"

	"os"
	"os/signal"
	"syscall"

	consul "github.com/hashicorp/consul/api"
)

func Register(name string, host string, port int, target string, interval time.Duration, ttl int) (<-chan struct{}, error) {
	conf := &consul.Config{Address: target}
	dereg := make(chan struct{})
	client, err := consul.NewClient(conf)
	if err != nil {
		return dereg, fmt.Errorf("can't create consul client: target = %s, err = %v", target, err)
	}
	serviceId := fmt.Sprintf("%s-%s-%d", name, host, port)
	stop := make(chan struct{})

	go func() {
		signalChan := make(chan os.Signal, 1)
		signal.Notify(signalChan, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL, syscall.SIGHUP, syscall.SIGQUIT)
		<-signalChan
		stop <- struct{}{}

		err := client.Agent().ServiceDeregister(serviceId)
		if err != nil {
			log.Printf("can't deregister service: serviceId = %s, err = %v", serviceId, err)
		}

		err = client.Agent().CheckDeregister(serviceId)
		if err != nil {
			log.Printf("can't check deregister: serviceId = %s, err = %v", serviceId, err)
		}
		dereg <- struct{}{}
	}()

	go func() {
		ticker := time.NewTicker(interval)
		for {
			select {
			case <-ticker.C:
				err = client.Agent().UpdateTTL(serviceId, "", "passing")
				if err != nil {
					log.Printf("can't update ttl: serviceId = %s, err = %v", serviceId, err)
				}
			case <-stop:
				return
			}
		}
	}()

	regs := &consul.AgentServiceRegistration{
		ID:      serviceId,
		Name:    name,
		Address: host,
		Port:    port,
	}
	if err = client.Agent().ServiceRegister(regs); err != nil {
		log.Printf("can't regster service: regs = %+v, err = %v", regs, err)
	}

	check := consul.AgentServiceCheck{
		TTL:    fmt.Sprintf("%ds", ttl),
		Status: "passing",
	}
	checkReg := &consul.AgentCheckRegistration{
		ID:                serviceId,
		Name:              name,
		ServiceID:         serviceId,
		AgentServiceCheck: check,
	}
	if err = client.Agent().CheckRegister(checkReg); err != nil {
		log.Printf("can't check register: regs = %+v, err = %v", checkReg, err)
	}
	return dereg, nil
}
