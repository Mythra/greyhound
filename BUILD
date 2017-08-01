load("@io_bazel_rules_go//go:def.bzl", "go_prefix")
go_prefix("github.com/instructure/dd-db-warden")
 
load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_binary", "go_test")
 
go_library(
  name = "greyhound-lib",
  srcs = glob([
    'src/**/*.go',
  ], exclude = [
    'src/**/*_test.go'
  ]),
  deps = [
    '@com_github_go_yaml_yaml//:go_default_library',
    '@com_github_syndtr_goleveldb//leveldb:go_default_library',
    '@com_github_spf13_afero//:go_default_library',
    '@com_github_cenkalti_backoff//:go_default_library',
  ],
  visibility = ["//visibility:public"]
)
 
go_binary(
  name = "greyhound",
  library = ':greyhound-lib',
  visibility = ["//visibility:public"],
)

go_test(
  name = "greyhound-tests",
  srcs = glob([
    'src/**/*_test.go'
  ]),
  deps = [
    "@com_github_h2non_gock//:go_default_library"
  ],
  library = ':greyhound-lib',
  size = "small"
)