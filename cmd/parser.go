// This file is part of go-pff (https://github.com/mooijtech/go-pff)
// Copyright (C) 2021 Marten Mooij (https://www.mooijtech.com/)
package main

import (
	log "github.com/sirupsen/logrus"
	pff "pff/pkg"
)

func main() {
	log.SetLevel(log.DebugLevel)
	log.Infof("Starting go-pff...")

	parser := pff.NewParser()

	parser.Parse("data/enron.pst")
}
