repos
-----

A simple cmd to show status of multiple git repositories.


* Create configuration listing your git repositories:

```
$ cat ~/.config/repos/repos
/home/fabio/go/src/github.com/fabiofalci/sconsify
/home/fabio/go/src/github.com/fabiofalci/gohit
/home/fabio/go/src/github.com/fabiofalci/i3what
/home/fabio/go/src/github.com/fabiofalci/money
/home/fabio/go/src/github.com/fabiofalci/repos
/home/fabio/go/src/github.com/fabiofalci/flagrc
/home/fabio/dev/p/agent
/home/fabio/dev/p/containers
/home/fabio/dev/p/cv
/home/fabio/dev/p/fabiofalci.github.io
/home/fabio/dev/p/practice
/home/fabio/dev/p/katas
/home/fabio/dev/p/test
```

And then run `repos`:

```
$ repos
                      Remot Local [branch]
             sconsify ----- CHANG [mock_libspotify]
                gohit ----- ----- [master]
               i3what ----- ----- [master]
                money ----- ----- [master]
                repos AHEAD UNTRA [master]
               flagrc ----- ----- [master]
                agent ----- ----- [master]
           containers ----- ----- [master]
                   cv ----- ----- [master]
 fabiofalci.github.io AHEAD UNTRA [jwt]
             practice ----- ----- [master]
                katas ----- ----- [master]
                 test NO-RE ----- [master]

```

Remote status can be: `AHEAD`, `BEHIND`, `NO-REMOTE` and `----` for in sync with remote.

Local status can be: `CHANGES`, `UNTRACKED` and `----` for no changes.

How to build
------------

Clone the repository, install `go` and

```
$ make build
```

Binary will be generated in the folder `bundles`.
