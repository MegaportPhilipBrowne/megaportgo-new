// Copyright 2020 Megaport Pty Ltd
//
// Licensed under the Mozilla Public License, Version 2.0 (the
// "License"); you may not use this file except in compliance with
// the License. You may obtain a copy of the License at
//
//       https://mozilla.org/MPL/2.0/
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package mega_err

const ERR_PORT_PROVISION_TIMEOUT_EXCEED = "the port took too long to provision"
const ERR_MCR_PROVISION_TIMEOUT_EXCEED = "the MCR took too long to provision"
const ERR_MVE_PROVISION_TIMEOUT_EXCEED = "the MVE took too long to provision"
const ERR_VXC_PROVISION_TIMEOUT_EXCEED = "the VXC took too long to provision"

const ERR_VXC_NOT_LIVE = "the VXC is not in the expected LIVE state"
const ERR_VXC_UPDATE_TIMEOUT_EXCEED = "the VXC took longer than 15 minutes to update, and has failed"
const ERR_WRONG_PRODUCT_MODIFY = "sorry you can only update Ports and MCR2 using this method"
const ERR_NO_AVAILABLE_VXC_PORTS = "there are no available ports for you to connect to"
const ERR_INVALID_PARTNER = "the partner type you have passed is not valid"
const ERR_TERM_NOT_VALID = "invalid term, valid values are 1, 12, 24, and 36"
const ERR_PORT_ALREADY_LOCKED = "that port is already locked, cannot lock"
const ERR_PORT_NOT_LOCKED = "that port not locked, cannot unlock"
const ERR_PORT_NOT_LIVE = "the port is not in the expected LIVE state"
const ERR_MCR_INVALID_PORT_SPEED = "invalid port speed, valid speeds are 1000, 2500, 5000, and 10000"
const ERR_MCR_NOT_LIVE = "the MCR is not in the expected LIVE state"
const ERR_LOCATION_NOT_FOUND = "unable to find location"
const ERR_NO_MATCHING_LOCATIONS = "unable to find location based on search criteria"
const ERR_NO_OTP_KEY_DEFINED string = "a one time password key is not defined and we cannot generate a OTP due to this"
const ERR_PARSING_ERR_RESPONSE = "status code '%v' received from api and there has been an error parsing an error: %s. " +
	"The error body was:\nBEGIN\n%v\nEND\n"
const ERR_PARTNER_PORT_NO_RESULTS = "sorry there were no results returned based on the given filters"
const ERR_SESSION_TOKEN_STILL_EXIST = "it looks like the session was not removed and still exists, logout did not work"
const ERR_MEGAPORT_URL_NOT_SET = "The variable megaport_url has not been set correctly"
