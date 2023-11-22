package core

import (
	"context"
	"fmt"
	"strings"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"mediagui/domain"
	"mediagui/lib"
	"mediagui/logger"
	pb "mediagui/mediaagent"
)

func (c *Core) PruneMovies() {
	t0 := time.Now()

	lib.Notify(c.ctx.Hub, "prune:begin", "Started Prune Process")

	options := &domain.Options{Offset: 0, Limit: 99999999999999, SortBy: "title", SortOrder: "asc"}
	_, items := c.storage.GetMovies(options)

	hostItems := make(map[string][]*pb.Item)

	for _, host := range c.ctx.UnraidHosts {
		hostItems[host] = make([]*pb.Item, 0)
	}

	for _, item := range items {
		index := strings.Index(item.Location, ":")
		if index == -1 {
			// a valid location is wopr:/mnt/user/films/bluray/22 Bullets (2010)/22.Bullets_BLURAY.iso
			// if a ':' isn't found, then this must be a stub
			continue
		}

		host := item.Location[:index]
		location := item.Location[index+1:]

		hostItems[host] = append(hostItems[host], &pb.Item{Id: item.ID, Location: location, Title: item.Title})
	}

	for _, host := range c.ctx.UnraidHosts {
		address := fmt.Sprintf("%s:7624", host)

		conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			logger.Yellow("Unable to connect to host (%s): %s", address, err)
			continue
		}
		defer conn.Close()

		client := pb.NewMediaAgentClient(conn)

		rsp, err := client.Exists(context.Background(), &pb.ExistsReq{Items: hostItems[host]})
		if err != nil {
			logger.Yellow("Unable to check exist (%s): %s", address, err)
			continue
		}

		for _, item := range rsp.Items {
			c.storage.DeleteMovie(item.Id)
			lib.Notify(c.ctx.Hub, "prune:delete", fmt.Sprintf("DELETED: [%d] %s (%s))", item.Id, item.Title, item.Location))
		}
	}

	lib.Notify(c.ctx.Hub, "prune:end", fmt.Sprintf("Prune process finished (%s elapsed)", time.Since(t0).String()))
}
