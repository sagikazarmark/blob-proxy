go_binary(
    name = "blob-proxy",
    srcs = glob(["*.go"], exclude = ["*_test.go"]),
    deps = [
        "//third_party/go:github.com__oklog__run",
        "//third_party/go:github.com__spf13__pflag",
        "//third_party/go:gocloud.dev__blob",
        "//third_party/go:gocloud.dev__blob__fileblob",
        "//third_party/go:gocloud.dev__gcerrors",
    ],
)
