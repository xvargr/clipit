package main

import (
	"log"

	"github.com/xvargr/clippit/internal/URLShortener"
	"github.com/xvargr/clippit/internal/config"
)

func PruneTask() {
	interval := config.GetConfig().PruneIntervalHour
	purged := URLShortener.Instance().Prune(interval)
	log.Default().Printf("Pruned expired mappings, purged %d entries\n", purged)
}
