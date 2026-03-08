package main

type Device struct {
	Alias string
	UUID  uuid
}

type Watchgroup struct {
	Name     string
	Overseer string
	Devices  []Device
}

type Notification struct {
	Header      string
	Description string
}
