/*
	Copyright 2019 Nokia

	Licensed under the Apache License, Version 2.0 (the "License");
	you may not use this file except in compliance with the License.
	You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

	Unless required by applicable law or agreed to in writing, software
	distributed under the License is distributed on an "AS IS" BASIS,
	WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
	See the License for the specific language governing permissions and
	limitations under the License.
*/

package govel

import "time"

// CollectorConfiguration parameters
type CollectorConfiguration struct {
	ServerRoot string `mapstructure:"serverRoot"`
	FQDN       string `mapstructure:"fqdn"`
	Port       int    `mapstructure:"port"`
	Secure     bool   `mapstructure:"secure"`
	Topic      string `mapstructure:"topic"`
	User       string `mapstructure:"user"`
	Password   string `mapstructure:"password"`
	PassPhrase string `mapsctructure:"passphrase,omitempty"` // passPhrase used to encrypt collector password in file
}

//NfcNamingCode mapping bettween NfcNamingCode (oam or etl) and Vnfcs
type NfcNamingCode struct {
	Type  string   `mapstructure:"type"`
	Vnfcs []string `mapstructure:"vnfcs"`
}

// EventConfiguration parameters
type EventConfiguration struct {
	VNFName             string          `mapstructure:"vnfName"`             // Name of this VNF, eg: dpa2bhsxp5001v
	ReportingEntityName string          `mapstructure:"reportingEntityName"` // Value of reporting entity field. Usually local VM (VNFC) name
	ReportingEntityID   string          `mapstructure:"reportingEntityID"`   // Value of reporting entity UUID. Usually local VM (VNFC) UUID
	MaxSize             int             `mapstructure:"maxSize,omitempty"`
	NfNamingCode        string          `mapstructure:"nfNamingCode,omitempty"` // "hspx"
	NfcNamingCodes      []NfcNamingCode `mapstructure:"nfcNamingCodes,omitempty"`
	RetryInterval       time.Duration   `mapstructure:"retryInterval,omitempty"`
	MaxMissed           int             `mapstructure:"maxMissed,omitempty"`
}
