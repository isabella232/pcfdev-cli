package downloader_test

import (
	"errors"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/golang/mock/gomock"
	dl "github.com/pivotal-cf/pcfdev-cli/downloader"
	"github.com/pivotal-cf/pcfdev-cli/downloader/mocks"
	"github.com/pivotal-cf/pcfdev-cli/pivnet"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Downloader", func() {
	var (
		downloader *dl.Downloader
		mockCtrl   *gomock.Controller
		mockClient *mocks.MockClient
		mockFS     *mocks.MockFS
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		mockClient = mocks.NewMockClient(mockCtrl)
		mockFS = mocks.NewMockFS(mockCtrl)

		downloader = &dl.Downloader{
			PivnetClient: mockClient,
			FS:           mockFS,
			ExpectedMD5:  "some-md5",
		}
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	Describe("#Download", func() {
		Context("when the file exists", func() {
			It("should not re-download the file", func() {
				gomock.InOrder(
					mockFS.EXPECT().CreateDir(filepath.Join("some-path")).Return(nil),
					mockFS.EXPECT().Exists(filepath.Join("some-path", "some-file.ova")).Return(true, nil),
				)

				Expect(downloader.Download(filepath.Join("some-path", "some-file.ova"))).To(Succeed())
			})
		})

		Context("when file and partial file do not exist", func() {
			It("should download the file", func() {
				readCloser := &pivnet.DownloadReader{ReadCloser: ioutil.NopCloser(strings.NewReader("some-ova-contents"))}
				gomock.InOrder(
					mockFS.EXPECT().CreateDir(filepath.Join("some-path")).Return(nil),
					mockFS.EXPECT().Exists(filepath.Join("some-path", "some-file.ova")).Return(false, nil),
					mockFS.EXPECT().Exists(filepath.Join("some-path", "some-file.ova.partial")).Return(false, nil),
					mockClient.EXPECT().DownloadOVA(int64(0)).Return(readCloser, nil),
					mockFS.EXPECT().Write(filepath.Join("some-path", "some-file.ova.partial"), readCloser).Return(nil),
					mockFS.EXPECT().MD5(filepath.Join("some-path", "some-file.ova.partial")).Return("some-md5", nil),
					mockFS.EXPECT().Move(filepath.Join("some-path", "some-file.ova.partial"), filepath.Join("some-path", "some-file.ova")),
				)

				Expect(downloader.Download(filepath.Join("some-path", "some-file.ova"))).To(Succeed())
			})
		})

		Context("when partial file does exist", func() {
			It("should resume the download of the partial file", func() {
				readCloser := &pivnet.DownloadReader{ReadCloser: ioutil.NopCloser(strings.NewReader("some-ova-contents"))}
				gomock.InOrder(
					mockFS.EXPECT().CreateDir(filepath.Join("some-path")).Return(nil),
					mockFS.EXPECT().Exists(filepath.Join("some-path", "some-file.ova")).Return(false, nil),
					mockFS.EXPECT().Exists(filepath.Join("some-path", "some-file.ova.partial")).Return(true, nil),
					mockFS.EXPECT().Length(filepath.Join("some-path", "some-file.ova.partial")).Return(int64(25), nil),
					mockClient.EXPECT().DownloadOVA(int64(25)).Return(readCloser, nil),
					mockFS.EXPECT().Write(filepath.Join("some-path", "some-file.ova.partial"), readCloser).Return(nil),
					mockFS.EXPECT().MD5(filepath.Join("some-path", "some-file.ova.partial")).Return("some-md5", nil),
					mockFS.EXPECT().Move(filepath.Join("some-path", "some-file.ova.partial"), filepath.Join("some-path", "some-file.ova")),
				)

				Expect(downloader.Download(filepath.Join("some-path", "some-file.ova"))).To(Succeed())
			})
		})

		Context("when partial file is downloaded but the checksum is not valid and the re-download succeeds", func() {
			It("should move the file to the downloaded path", func() {
				readCloser := &pivnet.DownloadReader{ReadCloser: ioutil.NopCloser(strings.NewReader("some-ova-contents"))}
				gomock.InOrder(
					mockFS.EXPECT().CreateDir(filepath.Join("some-path")).Return(nil),
					mockFS.EXPECT().Exists(filepath.Join("some-path", "some-file.ova")).Return(false, nil),
					mockFS.EXPECT().Exists(filepath.Join("some-path", "some-file.ova.partial")).Return(true, nil),
					mockFS.EXPECT().Length(filepath.Join("some-path", "some-file.ova.partial")).Return(int64(25), nil),
					mockClient.EXPECT().DownloadOVA(int64(25)).Return(readCloser, nil),
					mockFS.EXPECT().Write(filepath.Join("some-path", "some-file.ova.partial"), readCloser).Return(nil),
					mockFS.EXPECT().MD5(filepath.Join("some-path", "some-file.ova.partial")).Return("some-bad-md5", nil),
					mockFS.EXPECT().RemoveFile(filepath.Join("some-path", "some-file.ova.partial")).Return(nil),

					mockClient.EXPECT().DownloadOVA(int64(0)).Return(readCloser, nil),
					mockFS.EXPECT().Write(filepath.Join("some-path", "some-file.ova.partial"), readCloser).Return(nil),
					mockFS.EXPECT().MD5(filepath.Join("some-path", "some-file.ova.partial")).Return("some-md5", nil),
					mockFS.EXPECT().Move(filepath.Join("some-path", "some-file.ova.partial"), filepath.Join("some-path", "some-file.ova")),
				)

				Expect(downloader.Download(filepath.Join("some-path", "some-file.ova"))).To(Succeed())
			})
		})

		Context("when partial file is downloaded but the checksum is not valid and the re-download fails", func() {
			It("should return an error", func() {
				readCloser := &pivnet.DownloadReader{ReadCloser: ioutil.NopCloser(strings.NewReader("some-ova-contents"))}
				gomock.InOrder(
					mockFS.EXPECT().CreateDir(filepath.Join("some-path")).Return(nil),
					mockFS.EXPECT().Exists(filepath.Join("some-path", "some-file.ova")).Return(false, nil),
					mockFS.EXPECT().Exists(filepath.Join("some-path", "some-file.ova.partial")).Return(true, nil),
					mockFS.EXPECT().Length(filepath.Join("some-path", "some-file.ova.partial")).Return(int64(25), nil),
					mockClient.EXPECT().DownloadOVA(int64(25)).Return(readCloser, nil),
					mockFS.EXPECT().Write(filepath.Join("some-path", "some-file.ova.partial"), readCloser).Return(nil),
					mockFS.EXPECT().MD5(filepath.Join("some-path", "some-file.ova.partial")).Return("some-bad-md5", nil),
					mockFS.EXPECT().RemoveFile(filepath.Join("some-path", "some-file.ova.partial")).Return(nil),

					mockClient.EXPECT().DownloadOVA(int64(0)).Return(readCloser, nil),
					mockFS.EXPECT().Write(filepath.Join("some-path", "some-file.ova.partial"), readCloser).Return(nil),
					mockFS.EXPECT().MD5(filepath.Join("some-path", "some-file.ova.partial")).Return("some-bad-md5", nil),
				)

				Expect(downloader.Download(filepath.Join("some-path", "some-file.ova"))).To(MatchError("download failed"))
			})
		})

		Context("when creating the directory fails", func() {
			It("should return an error", func() {
				mockFS.EXPECT().CreateDir(filepath.Join("some-path")).Return(errors.New("some-error"))

				Expect(downloader.Download(filepath.Join("some-path", "some-file.ova"))).To(MatchError("some-error"))
			})
		})

		Context("when checking if the file exists", func() {
			It("should return an error", func() {
				gomock.InOrder(
					mockFS.EXPECT().CreateDir(filepath.Join("some-path")).Return(nil),
					mockFS.EXPECT().Exists(filepath.Join("some-path", "some-file.ova")).Return(false, errors.New("some-error")),
				)

				Expect(downloader.Download(filepath.Join("some-path", "some-file.ova"))).To(MatchError("some-error"))
			})
		})

		Context("when checking if the partial file exists", func() {
			It("should return an error", func() {
				gomock.InOrder(
					mockFS.EXPECT().CreateDir(filepath.Join("some-path")).Return(nil),
					mockFS.EXPECT().Exists(filepath.Join("some-path", "some-file.ova")).Return(false, nil),
					mockFS.EXPECT().Exists(filepath.Join("some-path", "some-file.ova.partial")).Return(false, errors.New("some-error")),
				)

				Expect(downloader.Download(filepath.Join("some-path", "some-file.ova"))).To(MatchError("some-error"))
			})
		})

		Context("when checking the length of the partial file fails", func() {
			It("should return an error", func() {
				gomock.InOrder(
					mockFS.EXPECT().CreateDir(filepath.Join("some-path")).Return(nil),
					mockFS.EXPECT().Exists(filepath.Join("some-path", "some-file.ova")).Return(false, nil),
					mockFS.EXPECT().Exists(filepath.Join("some-path", "some-file.ova.partial")).Return(true, nil),
					mockFS.EXPECT().Length(filepath.Join("some-path", "some-file.ova.partial")).Return(int64(0), errors.New("some-error")),
				)

				Expect(downloader.Download(filepath.Join("some-path", "some-file.ova"))).To(MatchError("some-error"))
			})
		})

		Context("when downloading the file fails", func() {
			It("should return an error", func() {
				gomock.InOrder(
					mockFS.EXPECT().CreateDir(filepath.Join("some-path")).Return(nil),
					mockFS.EXPECT().Exists(filepath.Join("some-path", "some-file.ova")).Return(false, nil),
					mockFS.EXPECT().Exists(filepath.Join("some-path", "some-file.ova.partial")).Return(false, nil),
					mockClient.EXPECT().DownloadOVA(int64(0)).Return(nil, errors.New("some-error")),
				)

				Expect(downloader.Download(filepath.Join("some-path", "some-file.ova"))).To(MatchError("some-error"))
			})
		})

		Context("when writing the downloaded file fails", func() {
			It("should return an error", func() {
				readCloser := &pivnet.DownloadReader{ReadCloser: ioutil.NopCloser(strings.NewReader("some-ova-contents"))}
				gomock.InOrder(
					mockFS.EXPECT().CreateDir(filepath.Join("some-path")).Return(nil),
					mockFS.EXPECT().Exists(filepath.Join("some-path", "some-file.ova")).Return(false, nil),
					mockFS.EXPECT().Exists(filepath.Join("some-path", "some-file.ova.partial")).Return(false, nil),
					mockClient.EXPECT().DownloadOVA(int64(0)).Return(readCloser, nil),
					mockFS.EXPECT().Write(filepath.Join("some-path", "some-file.ova.partial"), readCloser).Return(errors.New("some-error")),
				)

				Expect(downloader.Download(filepath.Join("some-path", "some-file.ova"))).To(MatchError("some-error"))
			})
		})

		Context("when checking the MD5 of the downloaded file fails", func() {
			It("should return an error", func() {
				readCloser := &pivnet.DownloadReader{ReadCloser: ioutil.NopCloser(strings.NewReader("some-ova-contents"))}
				gomock.InOrder(
					mockFS.EXPECT().CreateDir(filepath.Join("some-path")).Return(nil),
					mockFS.EXPECT().Exists(filepath.Join("some-path", "some-file.ova")).Return(false, nil),
					mockFS.EXPECT().Exists(filepath.Join("some-path", "some-file.ova.partial")).Return(false, nil),
					mockClient.EXPECT().DownloadOVA(int64(0)).Return(readCloser, nil),
					mockFS.EXPECT().Write(filepath.Join("some-path", "some-file.ova.partial"), readCloser).Return(nil),
					mockFS.EXPECT().MD5(filepath.Join("some-path", "some-file.ova.partial")).Return("", errors.New("some-error")),
				)

				Expect(downloader.Download(filepath.Join("some-path", "some-file.ova"))).To(MatchError("some-error"))
			})
		})

		Context("when removing the partial file fails", func() {
			It("should return an error", func() {
				readCloser := &pivnet.DownloadReader{ReadCloser: ioutil.NopCloser(strings.NewReader("some-ova-contents"))}
				gomock.InOrder(
					mockFS.EXPECT().CreateDir(filepath.Join("some-path")).Return(nil),
					mockFS.EXPECT().Exists(filepath.Join("some-path", "some-file.ova")).Return(false, nil),
					mockFS.EXPECT().Exists(filepath.Join("some-path", "some-file.ova.partial")).Return(true, nil),
					mockFS.EXPECT().Length(filepath.Join("some-path", "some-file.ova.partial")).Return(int64(25), nil),
					mockClient.EXPECT().DownloadOVA(int64(25)).Return(readCloser, nil),
					mockFS.EXPECT().Write(filepath.Join("some-path", "some-file.ova.partial"), readCloser).Return(nil),
					mockFS.EXPECT().MD5(filepath.Join("some-path", "some-file.ova.partial")).Return("some-bad-md5", nil),
					mockFS.EXPECT().RemoveFile(filepath.Join("some-path", "some-file.ova.partial")).Return(errors.New("some-error")),
				)

				Expect(downloader.Download(filepath.Join("some-path", "some-file.ova"))).To(MatchError("some-error"))
			})
		})

		Context("when the md5 of a file download does not match the expected md5", func() {
			It("should return an error", func() {
				readCloser := &pivnet.DownloadReader{ReadCloser: ioutil.NopCloser(strings.NewReader("some-ova-contents"))}
				gomock.InOrder(
					mockFS.EXPECT().CreateDir(filepath.Join("some-path")).Return(nil),
					mockFS.EXPECT().Exists(filepath.Join("some-path", "some-file.ova")).Return(false, nil),
					mockFS.EXPECT().Exists(filepath.Join("some-path", "some-file.ova.partial")).Return(false, nil),
					mockClient.EXPECT().DownloadOVA(int64(0)).Return(readCloser, nil),
					mockFS.EXPECT().Write(filepath.Join("some-path", "some-file.ova.partial"), readCloser).Return(nil),
					mockFS.EXPECT().MD5(filepath.Join("some-path", "some-file.ova.partial")).Return("some-bad-md5", nil),
				)

				Expect(downloader.Download(filepath.Join("some-path", "some-file.ova"))).To(MatchError("download failed"))
			})
		})

		Context("when downloading the file fails after downloading the partial file failed", func() {
			It("should return an error", func() {
				readCloser := &pivnet.DownloadReader{ReadCloser: ioutil.NopCloser(strings.NewReader("some-ova-contents"))}
				gomock.InOrder(
					mockFS.EXPECT().CreateDir(filepath.Join("some-path")).Return(nil),
					mockFS.EXPECT().Exists(filepath.Join("some-path", "some-file.ova")).Return(false, nil),
					mockFS.EXPECT().Exists(filepath.Join("some-path", "some-file.ova.partial")).Return(true, nil),
					mockFS.EXPECT().Length(filepath.Join("some-path", "some-file.ova.partial")).Return(int64(25), nil),
					mockClient.EXPECT().DownloadOVA(int64(25)).Return(readCloser, nil),
					mockFS.EXPECT().Write(filepath.Join("some-path", "some-file.ova.partial"), readCloser).Return(nil),
					mockFS.EXPECT().MD5(filepath.Join("some-path", "some-file.ova.partial")).Return("some-bad-md5", nil),
					mockFS.EXPECT().RemoveFile(filepath.Join("some-path", "some-file.ova.partial")).Return(nil),
					mockClient.EXPECT().DownloadOVA(int64(0)).Return(nil, errors.New("some-error")),
				)

				Expect(downloader.Download(filepath.Join("some-path", "some-file.ova"))).To(MatchError("some-error"))
			})
		})

		Context("when writing the file fails after downloading the partial file failed", func() {
			It("should return an error", func() {
				readCloser := &pivnet.DownloadReader{ReadCloser: ioutil.NopCloser(strings.NewReader("some-ova-contents"))}
				gomock.InOrder(
					mockFS.EXPECT().CreateDir(filepath.Join("some-path")).Return(nil),
					mockFS.EXPECT().Exists(filepath.Join("some-path", "some-file.ova")).Return(false, nil),
					mockFS.EXPECT().Exists(filepath.Join("some-path", "some-file.ova.partial")).Return(true, nil),
					mockFS.EXPECT().Length(filepath.Join("some-path", "some-file.ova.partial")).Return(int64(25), nil),
					mockClient.EXPECT().DownloadOVA(int64(25)).Return(readCloser, nil),
					mockFS.EXPECT().Write(filepath.Join("some-path", "some-file.ova.partial"), readCloser).Return(nil),
					mockFS.EXPECT().MD5(filepath.Join("some-path", "some-file.ova.partial")).Return("some-bad-md5", nil),
					mockFS.EXPECT().RemoveFile(filepath.Join("some-path", "some-file.ova.partial")).Return(nil),
					mockClient.EXPECT().DownloadOVA(int64(0)).Return(readCloser, nil),
					mockFS.EXPECT().Write(filepath.Join("some-path", "some-file.ova.partial"), readCloser).Return(errors.New("some-error")),
				)

				Expect(downloader.Download(filepath.Join("some-path", "some-file.ova"))).To(MatchError("some-error"))
			})
		})

		Context("when checking the MD5 of the file fails after downloading the partial file failed", func() {
			It("should return an error", func() {
				readCloser := &pivnet.DownloadReader{ReadCloser: ioutil.NopCloser(strings.NewReader("some-ova-contents"))}
				gomock.InOrder(
					mockFS.EXPECT().CreateDir(filepath.Join("some-path")).Return(nil),
					mockFS.EXPECT().Exists(filepath.Join("some-path", "some-file.ova")).Return(false, nil),
					mockFS.EXPECT().Exists(filepath.Join("some-path", "some-file.ova.partial")).Return(true, nil),
					mockFS.EXPECT().Length(filepath.Join("some-path", "some-file.ova.partial")).Return(int64(25), nil),
					mockClient.EXPECT().DownloadOVA(int64(25)).Return(readCloser, nil),
					mockFS.EXPECT().Write(filepath.Join("some-path", "some-file.ova.partial"), readCloser).Return(nil),
					mockFS.EXPECT().MD5(filepath.Join("some-path", "some-file.ova.partial")).Return("some-bad-md5", nil),
					mockFS.EXPECT().RemoveFile(filepath.Join("some-path", "some-file.ova.partial")).Return(nil),
					mockClient.EXPECT().DownloadOVA(int64(0)).Return(readCloser, nil),
					mockFS.EXPECT().Write(filepath.Join("some-path", "some-file.ova.partial"), readCloser).Return(nil),
					mockFS.EXPECT().MD5(filepath.Join("some-path", "some-file.ova.partial")).Return("", errors.New("some-error")),
				)

				Expect(downloader.Download(filepath.Join("some-path", "some-file.ova"))).To(MatchError("some-error"))
			})
		})

		Context("when moving the file fails", func() {
			It("should return an error", func() {
				readCloser := &pivnet.DownloadReader{ReadCloser: ioutil.NopCloser(strings.NewReader("some-ova-contents"))}
				gomock.InOrder(
					mockFS.EXPECT().CreateDir(filepath.Join("some-path")).Return(nil),
					mockFS.EXPECT().Exists(filepath.Join("some-path", "some-file.ova")).Return(false, nil),
					mockFS.EXPECT().Exists(filepath.Join("some-path", "some-file.ova.partial")).Return(false, nil),
					mockClient.EXPECT().DownloadOVA(int64(0)).Return(readCloser, nil),
					mockFS.EXPECT().Write(filepath.Join("some-path", "some-file.ova.partial"), readCloser).Return(nil),
					mockFS.EXPECT().MD5(filepath.Join("some-path", "some-file.ova.partial")).Return("some-md5", nil),
					mockFS.EXPECT().Move(filepath.Join("some-path", "some-file.ova.partial"), filepath.Join("some-path", "some-file.ova")).Return(errors.New("some-error")),
				)

				Expect(downloader.Download(filepath.Join("some-path", "some-file.ova"))).To(MatchError("some-error"))
			})
		})
	})
})