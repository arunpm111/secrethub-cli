package secrethub

import (
	"testing"
	"time"

	"github.com/secrethub/secrethub-cli/internals/cli/ui/fakeui"

	"github.com/secrethub/secrethub-go/internals/api"
	"github.com/secrethub/secrethub-go/internals/assert"
	"github.com/secrethub/secrethub-go/internals/errio"
	"github.com/secrethub/secrethub-go/pkg/secrethub"
	"github.com/secrethub/secrethub-go/pkg/secrethub/fakeclient"
)

func TestCredentialListCommand_Run(t *testing.T) {
	testErr := errio.Namespace("test").Code("test").Error("test error")

	cases := map[string]struct {
		cmd               CredentialListCommand
		credentialService fakeclient.CredentialService
		newClientErr      error
		expectedOut       string
		expectedErr       error
	}{
		"success": {
			cmd: CredentialListCommand{},
			credentialService: fakeclient.CredentialService{
				ListFunc: func(_ *secrethub.CredentialListParams) secrethub.CredentialIterator {
					return &fakeclient.CredentialIterator{
						Credentials: []*api.Credential{
							{
								Type:        "test",
								CreatedAt:   time.Now(),
								Fingerprint: "8E146D837D4CA1DC4315167B11A39C92",
								Description: "credential 1",
							},
							{
								CreatedAt:   time.Now(),
								Fingerprint: "DFC3D1F0D9842F17425403D0A9474C36",
								Description: "credential 2",
								Enabled:     true,
							},
						},
						CurrentIndex: 0,
						Err:          nil,
					}
				},
			},
			expectedOut: "FINGERPRINT       TYPE  ENABLED  CREATED                 DESCRIPTION\n" +
				"8E146D837D4CA1DC  test  no       Less than a second ago  credential 1\n" +
				"DFC3D1F0D9842F17        yes      Less than a second ago  credential 2\n",
		},
		"new client error": {
			newClientErr: testErr,
			expectedErr:  testErr,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			// Setup
			testIO := fakeui.NewIO(t)
			tc.cmd.io = testIO

			tc.cmd.newClient = func() (secrethub.ClientInterface, error) {
				return fakeclient.Client{
					CredentialService: &tc.credentialService,
				}, tc.newClientErr
			}

			// Run
			err := tc.cmd.Run()

			// Assert
			assert.Equal(t, err, tc.expectedErr)
			assert.Equal(t, testIO.Out.String(), tc.expectedOut)
		})
	}
}
