package vm_test

import (
	"errors"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/pivotal-cf/pcfdev-cli/config"
	"github.com/pivotal-cf/pcfdev-cli/ssh"
	"github.com/pivotal-cf/pcfdev-cli/vm"
	"github.com/pivotal-cf/pcfdev-cli/vm/mocks"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Paused", func() {
	var (
		mockCtrl *gomock.Controller
		mockUI   *mocks.MockUI
		mockVBox *mocks.MockVBox
		mockSSH  *mocks.MockSSH
		mockFS   *mocks.MockFS
		pausedVM vm.Paused
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		mockUI = mocks.NewMockUI(mockCtrl)
		mockVBox = mocks.NewMockVBox(mockCtrl)
		mockSSH = mocks.NewMockSSH(mockCtrl)
		mockFS = mocks.NewMockFS(mockCtrl)

		pausedVM = vm.Paused{
			VMConfig: &config.VMConfig{
				Name:    "some-vm",
				Domain:  "some-domain",
				IP:      "some-ip",
				SSHPort: "some-port",
			},
			VBox:      mockVBox,
			UI:        mockUI,
			SSHClient: mockSSH,
			FS:        mockFS,

			Config: &config.Config{
				PrivateKeyPath: "some-private-key-path",
			},
		}
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	Describe("Suspend", func() {
		It("should say a message", func() {
			mockUI.EXPECT().Say("Your VM is suspended and system memory for the VM is still allocated. Resume and suspend to suspend pcfdev VM to the disk.")
			Expect(pausedVM.Suspend()).To(Succeed())
		})
	})

	Describe("Stop", func() {
		It("should say a message", func() {
			mockUI.EXPECT().Say("Your VM is currently suspended. You must resume your VM with `cf dev resume` to shut it down.")
			Expect(pausedVM.Stop()).To(Succeed())
		})
	})

	Describe("Start", func() {
		It("should start vm", func() {
			addresses := []ssh.SSHAddress{
				{IP: "127.0.0.1", Port: "some-port"},
				{IP: "some-ip", Port: "22"},
			}
			gomock.InOrder(
				mockUI.EXPECT().Say("Resuming VM..."),
				mockVBox.EXPECT().ResumePausedVM(pausedVM.VMConfig),
				mockFS.EXPECT().Read("some-private-key-path").Return([]byte("some-private-key"), nil),
				mockSSH.EXPECT().WaitForSSH(addresses, []byte("some-private-key"), 5*time.Minute),
				mockUI.EXPECT().Say("PCF Dev is now running."),
			)

			Expect(pausedVM.Start(&vm.StartOpts{})).To(Succeed())
		})

		Context("when starting the vm fails", func() {
			It("should return an error", func() {
				gomock.InOrder(
					mockUI.EXPECT().Say("Resuming VM..."),
					mockVBox.EXPECT().ResumePausedVM(pausedVM.VMConfig).Return(errors.New("some-error")),
				)

				Expect(pausedVM.Start(&vm.StartOpts{})).To(MatchError("failed to resume VM: some-error"))
			})
		})

		Context("when reading the private key fails", func() {
			It("should return an error", func() {
				gomock.InOrder(
					mockUI.EXPECT().Say("Resuming VM..."),
					mockVBox.EXPECT().ResumePausedVM(pausedVM.VMConfig),
					mockFS.EXPECT().Read("some-private-key-path").Return(nil, errors.New("some-error")),
				)

				Expect(pausedVM.Start(&vm.StartOpts{})).To(MatchError("failed to resume VM: some-error"))
			})
		})

		Context("when waiting for SSH fails", func() {
			It("should return an error", func() {
				addresses := []ssh.SSHAddress{
					{IP: "127.0.0.1", Port: "some-port"},
					{IP: "some-ip", Port: "22"},
				}
				gomock.InOrder(
					mockUI.EXPECT().Say("Resuming VM..."),
					mockVBox.EXPECT().ResumePausedVM(pausedVM.VMConfig),
					mockFS.EXPECT().Read("some-private-key-path").Return([]byte("some-private-key"), nil),
					mockSSH.EXPECT().WaitForSSH(addresses, []byte("some-private-key"), 5*time.Minute).Return(errors.New("some-error")),
				)

				Expect(pausedVM.Start(&vm.StartOpts{})).To(MatchError("failed to resume VM: some-error"))
			})
		})
	})

	Describe("VerifyStartOpts", func() {
		Context("when desired memory is passed", func() {
			It("should return an error", func() {
				Expect(pausedVM.VerifyStartOpts(&vm.StartOpts{
					Memory: 4000,
				})).To(MatchError("memory cannot be changed once the vm has been created"))
			})
		})

		Context("when desired cores is passed", func() {
			It("should return an error", func() {
				Expect(pausedVM.VerifyStartOpts(&vm.StartOpts{
					CPUs: 2,
				})).To(MatchError("cores cannot be changed once the vm has been created"))
			})
		})

		Context("when desired services is passed", func() {
			It("should return an error", func() {
				Expect(pausedVM.VerifyStartOpts(&vm.StartOpts{
					Services: "redis",
				})).To(MatchError("services cannot be changed once the vm has been created"))
			})
		})

		Context("when no opts are passed", func() {
			It("should succeed", func() {
				Expect(pausedVM.VerifyStartOpts(&vm.StartOpts{})).To(Succeed())
			})
		})

		Context("when desired IP is passed", func() {
			It("should return an error", func() {
				Expect(pausedVM.VerifyStartOpts(&vm.StartOpts{
					IP: "some-ip",
				})).To(MatchError("the -i flag cannot be used if the VM has already been created"))
			})
		})

		Context("when desired domain is passed", func() {
			It("should return an error", func() {
				Expect(pausedVM.VerifyStartOpts(&vm.StartOpts{
					Domain: "some-domain",
				})).To(MatchError("the -d flag cannot be used if the VM has already been created"))
			})
		})
	})

	Describe("Resume", func() {
		It("should resume vm", func() {
			addresses := []ssh.SSHAddress{
				{IP: "127.0.0.1", Port: "some-port"},
				{IP: "some-ip", Port: "22"},
			}
			gomock.InOrder(
				mockUI.EXPECT().Say("Resuming VM..."),
				mockVBox.EXPECT().ResumePausedVM(pausedVM.VMConfig),
				mockFS.EXPECT().Read("some-private-key-path").Return([]byte("some-private-key"), nil),
				mockSSH.EXPECT().WaitForSSH(addresses, []byte("some-private-key"), 5*time.Minute),
				mockUI.EXPECT().Say("PCF Dev is now running."),
			)

			Expect(pausedVM.Resume()).To(Succeed())
		})

		Context("when waiting for SSH fails", func() {
			It("should return an error", func() {
				addresses := []ssh.SSHAddress{
					{IP: "127.0.0.1", Port: "some-port"},
					{IP: "some-ip", Port: "22"},
				}
				gomock.InOrder(
					mockUI.EXPECT().Say("Resuming VM..."),
					mockVBox.EXPECT().ResumePausedVM(pausedVM.VMConfig),
					mockFS.EXPECT().Read("some-private-key-path").Return([]byte("some-private-key"), nil),
					mockSSH.EXPECT().WaitForSSH(addresses, []byte("some-private-key"), 5*time.Minute).Return(errors.New("some-error")),
				)

				Expect(pausedVM.Resume()).To(MatchError("failed to resume VM: some-error"))
			})
		})

		Context("when reading the private key fails", func() {
			It("should return an error", func() {
				gomock.InOrder(
					mockUI.EXPECT().Say("Resuming VM..."),
					mockVBox.EXPECT().ResumePausedVM(pausedVM.VMConfig),
					mockFS.EXPECT().Read("some-private-key-path").Return(nil, errors.New("some-error")),
				)

				Expect(pausedVM.Resume()).To(MatchError("failed to resume VM: some-error"))
			})
		})

		Context("when starting the vm fails", func() {
			It("should return an error", func() {
				gomock.InOrder(
					mockUI.EXPECT().Say("Resuming VM..."),
					mockVBox.EXPECT().ResumePausedVM(pausedVM.VMConfig).Return(errors.New("some-error")),
				)

				Expect(pausedVM.Resume()).To(MatchError("failed to resume VM: some-error"))
			})
		})
	})

	Describe("Status", func() {
		It("should return 'Suspended' with an explanation", func() {
			Expect(pausedVM.Status()).To(Equal("Suspended - system memory for the VM is still allocated. Resume and suspend to suspend pcfdev VM to the disk."))
		})
	})

	Describe("GetDebugLogs", func() {
		It("should say a message", func() {
			mockUI.EXPECT().Say("Your VM is suspended. Resume to retrieve debug logs.")
			Expect(pausedVM.GetDebugLogs()).To(Succeed())
		})
	})

	Describe("Trust", func() {
		It("should say a message", func() {
			mockUI.EXPECT().Say("Your VM is suspended. Resume to trust VM certificates.")
			Expect(pausedVM.Trust(&vm.StartOpts{})).To(Succeed())
		})
	})

	Describe("Target", func() {
		It("should say a message", func() {
			mockUI.EXPECT().Say("Your VM is suspended. Resume to target PCF Dev.")
			Expect(pausedVM.Target(false)).To(Succeed())
		})
	})

	Describe("SSH", func() {
		It("should say a message", func() {
			mockUI.EXPECT().Say("Your VM is suspended. Resume to SSH to PCF Dev.")
			Expect(pausedVM.SSH()).To(Succeed())
		})
	})
})
