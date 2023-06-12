// Copyright (c) 2023 Tiago Melo. All rights reserved.
// Use of this source code is governed by the MIT License that can be found in
// the LICENSE file.
package screen

import (
	_ "embed"
	"fmt"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/pterm/pterm"
	"github.com/tiagomelo/realtime-data-kafka/stats"
)

// Banners.
var (
	//go:embed banners/kafkaconsumer.txt
	kafkaConsumerBanner string
	//go:embed banners/kafkaproducer.txt
	kafkaProducerBanner string
)

// For ease of unit testing.
var (
	ptermDefaultAreaStart    = pterm.DefaultArea.Start
	ptermDefaultCenterSprint = func(a ...interface{}) string {
		return pterm.DefaultCenter.Sprint(a...)
	}
	stopAreaPrinter = func(areaPrinter *pterm.AreaPrinter) error {
		return areaPrinter.Stop()
	}
	layoutSprint = func(layout *pterm.CenterPrinter, a ...interface{}) string {
		return layout.Sprint(a...)
	}
	updateAreaPrinter = func(areaPrinter *pterm.AreaPrinter, text ...interface{}) {
		areaPrinter.Update(text...)
	}
)

// Screen is an interface for managing screen content.
// It defines a contract for types that can update the content on a screen.
type Screen interface {
	UpdateContent(finalUpdate bool) error
}

// baseScreen represents the basic screen structure.
type baseScreen struct {
	areaPrinter *pterm.AreaPrinter
	layout      *pterm.CenterPrinter
}

// KafkaConsumerScreen implements the Screen interface for consumer-related content.
type KafkaConsumerScreen struct {
	*baseScreen
	stats *stats.KafkaConsumerStats
}

// KafkaProducerScreen implements the Screen interface for producer-related content.
type KafkaProducerScreen struct {
	*baseScreen
	stats *stats.KafkaProducerStats
}

// NewKafkaConsumerScreen creates a new KafkaConsumerScreen.
func NewKafkaConsumerScreen(stats *stats.KafkaConsumerStats) (Screen, error) {
	baseScreen, err := createBaseScreen()
	if err != nil {
		return nil, err
	}
	return &KafkaConsumerScreen{
		baseScreen: baseScreen,
		stats:      stats,
	}, nil
}

// UpdateContent updates the content of the KafkaConsumerScreen.
func (s *KafkaConsumerScreen) UpdateContent(finalUpdate bool) error {
	out := []string{
		template("Total processed transactions", fmt.Sprintf("%d", s.stats.TotalTransactions())),
		template("Suspicous transactions", fmt.Sprintf("%d", s.stats.TotalSuspiciousTransactions())),
		template("Invalid kafka messages", fmt.Sprintf("%d", s.stats.TotalUnmarshallingMsgErrors())),
		template("Total DB errors", fmt.Sprintf("%d", s.stats.TotalInsertSuspiciousTransactionErrors())),
		template("Elapsed Time", formatDuration(s.stats.ElapsedTime())),
	}
	banner := ptermDefaultCenterSprint(string(kafkaConsumerBanner))
	content := layoutSprint(s.layout, strings.Join(out, "\n"))
	updateAreaPrinter(s.areaPrinter, banner+content)
	if finalUpdate {
		if err := stopAreaPrinter(s.areaPrinter); err != nil {
			return errors.Wrap(err, "stopping printer")
		}
	}
	return nil
}

// NewKafkaProducerScreen creates a new KafkaProducerScreen.
func NewKafkaProducerScreen(stats *stats.KafkaProducerStats) (Screen, error) {
	baseScreen, err := createBaseScreen()
	if err != nil {
		return nil, err
	}
	return &KafkaProducerScreen{
		baseScreen: baseScreen,
		stats:      stats,
	}, nil
}

// UpdateContent updates the content of the KafkaProducerScreen.
func (s *KafkaProducerScreen) UpdateContent(finalUpdate bool) error {
	out := []string{
		template("Total published messages", fmt.Sprintf("%d", s.stats.TotalPublishedMessages())),
		template("Total message delivery errors", fmt.Sprintf("%d", s.stats.TotalFailedMessageDeliveries())),
		template("Elapsed Time", formatDuration(s.stats.ElapsedTime())),
	}
	banner := ptermDefaultCenterSprint(string(kafkaProducerBanner))
	content := layoutSprint(s.layout, strings.Join(out, "\n"))
	updateAreaPrinter(s.areaPrinter, banner+content)
	if finalUpdate {
		if err := stopAreaPrinter(s.areaPrinter); err != nil {
			return errors.Wrap(err, "stopping printer")
		}
	}
	return nil
}

// template is a helper function to define the layout.
func template(name string, value string) string {
	const MAX_LENGTH int = 42
	pad := MAX_LENGTH - len(name)

	// Example output
	// [ Total                                  1883 ]
	return fmt.Sprintf("[ %s %*s ]", name, pad, value)
}

// formatDuration formats a duration to the format "hh:mm:ss".
func formatDuration(d time.Duration) string {
	d = d.Round(time.Second)
	h := d / time.Hour
	d -= h * time.Hour
	m := d / time.Minute
	d -= m * time.Minute
	s := d / time.Second
	return fmt.Sprintf("%02dh%02dm%02ds", h, m, s)
}

// createBaseScreen creates a new baseScreen.
func createBaseScreen() (*baseScreen, error) {
	area, err := ptermDefaultAreaStart()
	if err != nil {
		return nil, errors.Wrap(err, "starting printer")
	}
	return &baseScreen{area, new(pterm.CenterPrinter)}, nil
}
