timestamp = git_show("%ct")

date_fmt = "+%FT%T%z"

go_binary(
    name = "blob-proxy",
    srcs = glob(["*.go"], exclude = ["*_test.go"]),
    definitions = {
        "main.version": "${VERSION:-" + git_branch() + "}",
        "main.commitHash": git_commit()[0:8],
        "main.buildDate": f'$(date -u -d "@{timestamp}" "{date_fmt}" 2>/dev/null || date -u -r "{timestamp}" "{date_fmt}" 2>/dev/null || date -u "{date_fmt}")',
    },
    deps = [
        "//third_party/go:github.com__oklog__run",
        "//third_party/go:github.com__spf13__pflag",
        "//third_party/go:gocloud.dev__blob",
        "//third_party/go:gocloud.dev__blob__fileblob",
        "//third_party/go:gocloud.dev__blob__s3blob",
        "//third_party/go:gocloud.dev__gcerrors",
    ],
)
