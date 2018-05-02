// Package discovery provides a discovery service for chromecast devices
package discovery

import (
	"context"
	"strings"
	"time"

	"github.com/HeavyHorst/go-cast"
	"github.com/HeavyHorst/go-cast/log"
	"github.com/grandcat/zeroconf"
)

type Service struct {
	found     chan *cast.Client
	entriesCh chan *zeroconf.ServiceEntry
}

func NewService(ctx context.Context) *Service {
	s := &Service{
		found:     make(chan *cast.Client),
		entriesCh: make(chan *zeroconf.ServiceEntry, 10),
	}

	go s.listener(ctx)
	return s
}

func wrapChan(c chan *zeroconf.ServiceEntry) chan *zeroconf.ServiceEntry {
	entries := make(chan *zeroconf.ServiceEntry)
	go func() {
		for entry := range entries {
			c <- entry
		}
	}()
	return entries
}

func (d *Service) Run(ctx context.Context, interval time.Duration) error {
	resolver, err := zeroconf.NewResolver(nil)
	if err != nil {
		log.Fatalln("Failed to initialize resolver:", err.Error())
		return err
	}

	err = resolver.Browse(ctx, "_googlecast._tcp", "local", wrapChan(d.entriesCh))
	if err != nil {
		log.Fatalln("Failed to browse:", err.Error())
	}

	ticker := time.NewTicker(interval)
	for {
		select {
		case <-ticker.C:
			resolver.Browse(ctx, "_googlecast._tcp", "local", wrapChan(d.entriesCh))
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (d *Service) Found() chan *cast.Client {
	return d.found
}

func (d *Service) listener(ctx context.Context) {
	for entry := range d.entriesCh {
		log.Printf("New entry: %#v\n", entry)
		client := cast.NewClient(entry.AddrIPv4[0], entry.Port)

		info := decodeTxtRecord(strings.Join(entry.Text, "|"))
		client.SetName(info["fn"])
		client.SetInfo(info)

		select {
		case d.found <- client:
		case <-time.After(time.Second):
		case <-ctx.Done():
			break
		}
	}
}

func decodeTxtRecord(txt string) map[string]string {
	m := make(map[string]string)

	s := strings.Split(txt, "|")
	for _, v := range s {
		s := strings.Split(v, "=")
		if len(s) == 2 {
			m[s[0]] = s[1]
		}
	}

	return m
}
