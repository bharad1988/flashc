package cmd

import (
	"github.com/bharad1988/flashc/common"
)

// structre to be used for Image command
var superimag common.SuperImage

// attributes related to image
var imagename string

// UUID variables
var (
	batchUUID     string
	agentUUID     string
	containerUUID string
)

// Snapshot variables
var destroySnaps string
var snapname string

var lCon common.SuperCon

var respStatus string

var httpResp common.HTTPRetObj
