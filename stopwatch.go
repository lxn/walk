// Copyright 2019 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build windows

package walk

import (
	"bytes"
	"fmt"
	"sort"
	"sync"
	"text/tabwriter"
	"time"
)

type stopwatchItem struct {
	stopwatchStats
	subject     string
	startedTime time.Time
}

type stopwatchStats struct {
	count int64
	min   time.Duration
	max   time.Duration
	total time.Duration
}

func (sws *stopwatchStats) Average() time.Duration {
	return time.Nanosecond * time.Duration(sws.total.Nanoseconds()/sws.count)
}

type Stopwatch struct {
	mutex        sync.Mutex
	subject2item map[string]*stopwatchItem
}

func NewStopwatch() *Stopwatch {
	return &Stopwatch{
		subject2item: make(map[string]*stopwatchItem),
	}
}

func (sw *Stopwatch) Start(subject string) time.Time {
	sw.mutex.Lock()
	defer sw.mutex.Unlock()

	item, ok := sw.subject2item[subject]
	if !ok {
		item = &stopwatchItem{subject: subject}
		sw.subject2item[subject] = item
	}

	item.startedTime = time.Now()

	return item.startedTime
}

func (sw *Stopwatch) Stop(subject string) time.Duration {
	sw.mutex.Lock()
	defer sw.mutex.Unlock()

	item, ok := sw.subject2item[subject]
	if !ok || item.startedTime.IsZero() {
		return 0
	}

	duration := time.Now().Sub(item.startedTime)

	item.count++
	if duration < item.min || item.min == 0 {
		item.min = duration
	}
	if duration > item.max {
		item.max = duration
	}
	item.total += duration
	item.startedTime = time.Time{}

	return duration
}

func (sw *Stopwatch) Cancel(subject string) {
	sw.mutex.Lock()
	defer sw.mutex.Unlock()

	item, ok := sw.subject2item[subject]
	if !ok {
		return
	}

	item.startedTime = time.Time{}
}

func (sw *Stopwatch) Clear() {
	sw.mutex.Lock()
	defer sw.mutex.Unlock()

	for key := range sw.subject2item {
		delete(sw.subject2item, key)
	}
}

func (sw *Stopwatch) Print() {
	sw.mutex.Lock()

	items := make([]*stopwatchItem, 0, len(sw.subject2item))
	for _, item := range sw.subject2item {
		items = append(items, item)
	}

	sw.mutex.Unlock()

	sort.Slice(items, func(i, j int) bool {
		return items[i].total > items[j].total
	})

	var buf bytes.Buffer

	writer := tabwriter.NewWriter(&buf, 0, 8, 2, ' ', tabwriter.AlignRight)

	fmt.Fprintln(writer, "#\tSubject\tAverage\tTotal\tMin\tMax\t\tCount")

	for i, item := range items {
		fmt.Fprintf(writer, "%d\t%s\t%s\t%s\t%s\t%s\t\t%d\n", i+1, item.subject, item.Average(), item.total, item.min, item.max, item.count)
	}

	writer.Flush()

	fmt.Print(buf.String())
}
