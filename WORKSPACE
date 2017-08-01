git_repository(
  name = "io_bazel_rules_go",
  remote = "https://github.com/bazelbuild/rules_go.git",
  tag = "0.4.3",
)
load("@io_bazel_rules_go//go:def.bzl", "go_repositories", "new_go_repository")
 
go_repositories()

new_go_repository(
  name = "com_github_go_yaml_yaml",
  commit = "cd8b52f8269e0feb286dfeef29f8fe4d5b397e0b",
  importpath = "gopkg.in/yaml.v2"
)

new_go_repository(
  name = "com_github_syndtr_goleveldb",
  commit = "8c81ea47d4c41a385645e133e15510fc6a2a74b4",
  importpath = "github.com/syndtr/goleveldb"
)

new_go_repository(
  name = "com_github_spf13_afero",
  commit = "9be650865eab0c12963d8753212f4f9c66cdcf12",
  importpath = "github.com/spf13/afero"
)

new_go_repository(
  name = "com_github_h2non_gock",
  commit = "c67415ca9149fa7d71c70c22d6b42616fc282f59",
  importpath = "gopkg.in/h2non/gock.v1",
  vcs = "git",
  remote = "https://github.com/h2non/gock.git"
)

new_go_repository(
  name = "com_github_cenkalti_backoff",
  commit = "5d150e7eec023ce7a124856b37c68e54b4050ac7",
  importpath = "github.com/cenkalti/backoff"
)

# Goleveldb dependencies.

new_go_repository(
  name = "com_github_golang_snappy",
  commit = "553a641470496b2327abcac10b36396bd98e45c9",
  importpath = "github.com/golang/snappy"
)

# AFERO Dependencies

new_go_repository(
  name = "org_golang_x_text",
  importpath = "golang.org/x/text",
  remote = "https://github.com/golang/text.git",
  vcs = "git",
  commit = "19e51611da83d6be54ddafce4a4af510cb3e9ea4",
)

# Backoff Dependencies

new_go_repository(
  name = "org_golang_x_net",
  importpath = "golang.org/x/net",
  remote = "https://github.com/golang/net.git",
  vcs = "git",
  commit = "5b58a9c3e1690d33a592e5b791638e25eb9b3f70"
)