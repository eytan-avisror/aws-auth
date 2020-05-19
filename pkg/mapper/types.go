/*

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

package mapper

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"time"

	"k8s.io/client-go/kubernetes"
)

func init() {
	log.SetFlags(0)

}

type AuthMapper struct {
	KubernetesClient kubernetes.Interface
	LoggingEnabled   bool
}

func New(client kubernetes.Interface, isCommandline bool) *AuthMapper {
	var mapper = &AuthMapper{}
	mapper.KubernetesClient = client

	if !isCommandline {
		log.SetOutput(ioutil.Discard)
	}
	return mapper
}

var (
	DefaultRetryerBackoffFactor float64 = 2.0
	DefaultRetryerBackoffJitter         = true
)

// AwsAuthData represents the data of the aws-auth configmap
type AwsAuthData struct {
	MapRoles []*RolesAuthMap `yaml:"mapRoles"`
	MapUsers []*UsersAuthMap `yaml:"mapUsers"`
}

// SetMapRoles sets the MapRoles element
func (m *AwsAuthData) SetMapRoles(authMap []*RolesAuthMap) {
	m.MapRoles = authMap
}

// SetMapUsers sets the MapUsers element
func (m *AwsAuthData) SetMapUsers(authMap []*UsersAuthMap) {
	m.MapUsers = authMap
}

// RemoveArguments are the arguments for removing a mapRole or mapUsers
type RemoveArguments struct {
	KubeconfigPath string
	MapRoles       bool
	MapUsers       bool
	Username       string
	RoleARN        string
	UserARN        string
	Groups         []string
	WithRetries    bool
	MinRetryTime   time.Duration
	MaxRetryTime   time.Duration
	MaxRetryCount  int
	IsGlobal       bool
}

func (args *RemoveArguments) Validate() {
	if args.WithRetries {
		if args.MaxRetryCount < 1 {
			log.Fatal("error: --retry-max-count is invalid, must be greater than zero")
		}
	}

	if args.RoleARN == "" && args.MapRoles {
		log.Fatal("error: --rolearn not provided")
	}

	if args.UserARN == "" && args.MapUsers {
		log.Fatal("error: --userarn not provided")
	}

	if args.MapUsers && args.MapRoles {
		log.Fatal("error: --mapusers and --maproles are mutually exclusive")
	}

	if !args.MapUsers && !args.MapRoles {
		if !args.IsGlobal {
			log.Fatal("error: must select --mapusers or --maproles")
		}
	}
}

// UpsertArguments are the arguments for upserting a mapRole or mapUsers
type UpsertArguments struct {
	KubeconfigPath string
	MapRoles       bool
	MapUsers       bool
	Username       string
	RoleARN        string
	UserARN        string
	Groups         []string
	WithRetries    bool
	MinRetryTime   time.Duration
	MaxRetryTime   time.Duration
	MaxRetryCount  int
}

func (args *UpsertArguments) Validate() {
	if args.WithRetries {
		if args.MaxRetryCount < 1 {
			log.Fatal("error: --retry-max-count is invalid, must be greater than zero")
		}
	}

	if args.RoleARN == "" && args.MapRoles {
		log.Fatal("error: --rolearn not provided")
	}

	if args.UserARN == "" && args.MapUsers {
		log.Fatal("error: --userarn not provided")
	}

	if args.Username == "" {
		log.Fatal("error: --username not provided")
	}

	if args.MapUsers && args.MapRoles {
		log.Fatal("error: --mapusers and --maproles are mutually exclusive")
	}

	if !args.MapUsers && !args.MapRoles {
		log.Fatal("error: must select --mapusers or --maproles")
	}
}

// RolesAuthMap is the basic structure of a mapRoles authentication object
type RolesAuthMap struct {
	RoleARN  string   `yaml:"rolearn"`
	Username string   `yaml:"username"`
	Groups   []string `yaml:"groups,omitempty"`
}

func (r *RolesAuthMap) String() string {
	var s strings.Builder
	s.WriteString(fmt.Sprintf("- rolearn: %v\n  ", r.RoleARN))
	s.WriteString(fmt.Sprintf("username: %v\n  ", r.Username))
	s.WriteString("groups:\n")
	for _, group := range r.Groups {
		s.WriteString(fmt.Sprintf("  - %v\n", group))
	}
	return s.String()
}

// UsersAuthMap is the basic structure of a mapUsers authentication object
type UsersAuthMap struct {
	UserARN  string   `yaml:"userarn"`
	Username string   `yaml:"username"`
	Groups   []string `yaml:"groups,omitempty"`
}

func (r *UsersAuthMap) String() string {
	var s strings.Builder
	s.WriteString(fmt.Sprintf("- userarn: %v\n  ", r.UserARN))
	s.WriteString(fmt.Sprintf("username: %v\n  ", r.Username))
	s.WriteString("groups:\n")
	for _, group := range r.Groups {
		s.WriteString(fmt.Sprintf("  - %v\n", group))
	}
	return s.String()
}

// NewRolesAuthMap returns a new NewRolesAuthMap
func NewRolesAuthMap(rolearn, username string, groups []string) *RolesAuthMap {
	return &RolesAuthMap{
		RoleARN:  rolearn,
		Username: username,
		Groups:   groups,
	}
}

// NewUsersAuthMap returns a new NewUsersAuthMap
func NewUsersAuthMap(userarn, username string, groups []string) *UsersAuthMap {
	return &UsersAuthMap{
		UserARN:  userarn,
		Username: username,
		Groups:   groups,
	}
}

// SetUsername sets the Username value
func (r *UsersAuthMap) SetUsername(v string) *UsersAuthMap {
	r.Username = v
	return r
}

// SetGroups sets the Groups value
func (r *UsersAuthMap) SetGroups(g []string) *UsersAuthMap {
	r.Groups = g
	return r
}

// SetUsername sets the Username value
func (r *RolesAuthMap) SetUsername(v string) *RolesAuthMap {
	r.Username = v
	return r
}

// SetGroups sets the Groups value
func (r *RolesAuthMap) SetGroups(g []string) *RolesAuthMap {
	r.Groups = g
	return r
}
