package config

//This is for defaults as well as individual
//share polling, which override the defaults
//if present.
type pollDetails struct {
	pollStart    string
	pollEnd      string
	pollInterval int
}

//The share struct includes poll details for the
//individual shares (overrrides defaults)
type share struct {
	code        string
	description string
	sharePoll   pollDetails
}

type configDetails struct {
	pollDefaults pollDetails
	shares       []share
}

type config interface {
	readConfig()
}

type configFromDB struct {
	loadedConfig configDetails
}

type configFromFile struct {
	loadedCongig configDetails
}

func (c configFromDB) readConfig() {
	//Code to read from DB goes here.
}

func (c configFromFile) readConfig() {

}

//Up to here we have the individual configs loaded.
//Now need to create common routines
