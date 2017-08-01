# GreyHound #

GreyHound is a way to define in YAML (and thus in git) datadog dashboards. These are useful for setting up self
service dashboards that are automatically generated. Even for people that don't have access to Datadog. So you can
give say engineers read access, but not write and they can still go around building their horribly colored dashboards.

## Building Greyhound ##

Greyhound uses Bazel to build itself, so you should go over the Bazel Install Instructions: [HERE][BAZEL_INSTALL].
Once you've setup bazel you can come back here. You got Bazel? Good. The important thing to know about our bazel
integration is this won't mess with anything in your `$GOPATH`. Bazel builds/downloads dependencies in a completely
seperate area. As to give each build a reproducible effect while not screwing with anything in your system.

For linux (and I presume Mac), these files will be in: `~/.cache/bazel/<unique-project-id>/external/<projects are here>`.

Once you've installed bazel, and now about projects being stored in a different place you'll want to make sure that
you're running a semi up to date Go.

The Go versions this were tested with is:
v1.8.1

Once you've got an up to date Go, Installed Bazel, and know about it not interfering you can build the project with:

```
$ bazel build :greyhound
```

TADAH! You've now built our project. The first build might be kind of slow as it downloads/builds all of our dependencies,
but after that builds should be super quick. The final binary will be at: `$greyhound_path/bazel-bin/greyhound`. From
here you can take it upload it, and do whatever you want with it.

It should be noted you can also install [BazelD][bazled_link] to make things easier so you can just run:

```
$ bazled build
```

## Testing Greyhound ##

Testing is also provided by bazel, so make sure you've followed the instructions to install bazel as listed in the
building section of this guide. From there you can simply run:

```
$ bazel test :greyhound-tests
```

Then you'll have all of our tests automatically run.

It should also note this also benifits from [BazelD][bazled_link] to make things easier so you can just run:

```
$ bazled test
```

[BAZEL_INSTALL]: https://bazel.build/versions/master/docs/install.html
[bazled_link]: https://github.com/SecurityInsanity/bazled