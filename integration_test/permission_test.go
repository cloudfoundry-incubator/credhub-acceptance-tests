package integration_test

import (
	"encoding/json"
	"fmt"
	"time"

	. "github.com/cloudfoundry-incubator/credhub-acceptance-tests/test_helpers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gexec"
)


type Permission struct {
	Actor string `json: actor`
	Path string		`json: path`
	Operations []string `json: operations`
}

var _ = Describe("Permission Test", func() {
	var actor string
	var path string

	BeforeEach(func() {
		timestamp := fmt.Sprintf("%d", time.Now().UnixNano())
		actor = "actor-" + timestamp
		path = "/path-" + timestamp
	})

	Context("Set Permission", func() {
		Context("Set Permission is called on permission that does not exist", func() {
			It("creates a new permission", func() {
				session := RunCommand("set-permission", "-a", actor, "-p", path, "-o", "read, write")
				var permission Permission
				_ = json.Unmarshal(session.Out.Contents(), &permission)
				Expect(permission).Should(Equal(
					Permission{
						Actor: actor,
						Path: path,
						Operations: []string{"read", "write"},
					}))
				Eventually(session).Should(Exit(0))
			})
		})
		Context("Set Permission is called on permission that exists", func() {
			It("updates existing permission", func() {
				RunCommand("set-permission", "-a", actor, "-p", path, "-o", "read, write")
				session := RunCommand("set-permission", "-a", actor, "-p", path, "-o", "read, write, delete")
				var permission Permission
				_ = json.Unmarshal(session.Out.Contents(), &permission)
				Expect(permission).Should(Equal(
					Permission{
						Actor: actor,
						Path: path,
						Operations: []string{"read", "write", "delete"},
					}))
				Eventually(session).Should(Exit(0))
			})
		})
	})
	Context("Get Permission", func() {
		Context("Get permission called on permission that does not exist", func() {
			It("throws an error", func() {
				session := RunCommand("get-permission", "-a", actor, "-p", path)
				Eventually(session).Should(Exit(1))
			})
		})
		Context("Get permission called on permission that exists", func() {
			It("returns the permission", func() {
				RunCommand("set-permission", "-a", actor, "-p", path, "-o", "read, write")
				session := RunCommand("get-permission", "-a", actor, "-p", path)
				var permission Permission
				_ = json.Unmarshal(session.Out.Contents(), &permission)
				Expect(permission).Should(Equal(
					Permission{
						Actor: actor,
						Path: path,
						Operations: []string{"read", "write"},
					}))
				Eventually(session).Should(Exit(0))
			})
		})
	})
	Context("Delete Permission", func() {
		Context("Delete permission called on permission that does not exist", func() {
			It("throws an error", func() {
				session := RunCommand("delete-permission", "-a", actor, "-p", path)
				Eventually(session).Should(Exit(1))
			})
		})
		Context("Delete permission called on permission that exists", func() {
			It("returns the deleted permission", func() {
				RunCommand("set-permission", "-a", actor, "-p", path, "-o", "read, write")
				session := RunCommand("delete-permission", "-a", actor, "-p", path)
				var permission Permission
				_ = json.Unmarshal(session.Out.Contents(), &permission)
				Expect(permission).Should(Equal(
					Permission{
						Actor: actor,
						Path: path,
						Operations: []string{"read", "write"},
					}))
				Eventually(session).Should(Exit(0))
			})
		})
	})
})
