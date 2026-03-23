# The Fuck [![Build Status][workflow-badge]][workflow-link] [![MIT License][license-badge]](LICENSE.md)

*The Fuck* is a magnificent app, ported from nvbn's amazing
[python version](https://github.com/nvbn/thefuck)
which was inspired by a
[@liamosaur](https://web.archive.org/web/20200415095708/https://twitter.com/liamosaur/)
[tweet](https://web.archive.org/web/20150222032139/https://twitter.com/liamosaur/status/506975850596536320),
that corrects errors in previous console commands.  This mainly exists
because installing python packages keeps getting worse.  This makes
me sad, but it seems like the python community can't figure this out.
Or figures it out in multiple, incompatible ways.

[![gif with examples][examples-link]][examples-link]

More examples:

```bash
$ apt-get install vim
E: Could not open lock file /var/lib/dpkg/lock - open (13: Permission denied)
E: Unable to lock the administration directory (/var/lib/dpkg/), are you root?

$ fuck
sudo apt-get install vim [enter/ctrl-c]
[sudo] password for nvbn:
Reading package lists... Done
...
```

```bash
$ git push
fatal: The current branch master has no upstream branch.
To push the current branch and set the remote as upstream, use

    git push --set-upstream origin master

$ fuck
git push --set-upstream origin master [enter/ctrl-c]
Counting objects: 9, done.
...
```

```bash
$ git brnch
git: 'brnch' is not a git command. See 'git --help'.

Did you mean this?
    branch

$ fuck
git branch [enter/ctrl-c]
* master
```

## Contents

1. [Installation](#installation)
2. [Updating](#updating)
3. [How it works](#how-it-works)
4. [Contributing rules](#contributing-rules)
5. [Developing](#developing)
6. [License](#license-mit)

## Installation

Install with Go:

```bash
go install github.com/lyda/thefuck@main
```

Then add the alias to your shell config:

| Shell          | Files to edit                    | What to add                        |
|----------------|----------------------------------|------------------------------------|
| **bash**       | `~/.bashrc` or `~/.bash_profile` | `eval "$(thefuck)"`                |
| **zsh**        | `~/.zshrc`                       | `eval "$(thefuck)"`                |
| **fish**       | `~/.config/fish/config.fish`     | `thefuck \| source`                |
| **tcsh**       | `~/.tcshrc`                      | `eval $(thefuck)`                  |
| **PowerShell** | `$profile`                       | `iex "$(thefuck init powershell)"` |

> **Note for csh users and users of older tcsh shells**: You'll want
to replace `$()` with backticks.

> **Note for PowerShell users**: *The Fuck* will helpfully suggest a superior
> operating system. This is the full extent of PowerShell support and we
> consider it a feature.

Changes are only available in a new shell session. To apply immediately,
source your config file (e.g. `source ~/.bashrc`) or just type the
one-liner above.

##### [Back to Contents](#contents)

## Updating

```bash
go install github.com/lyda/thefuck@main
```

##### [Back to Contents](#contents)

## How it works

*The Fuck* re-runs the previous command and matches the output against a set of
rules. If a match is found, a new corrected command is created and presented for
confirmation. Press Enter to run it, or Ctrl-C to abort.

The following rules are included:

  * `adb_unknown_command` &ndash; fixes misspelled commands like `adb logcta`;
  * `ag_literal` &ndash; adds `-Q` to `ag` when suggested;
  * `apt_get` &ndash; installs app from apt if it is not installed;
  * `apt_get_search` &ndash; changes trying to search using `apt-get` with searching using `apt-cache`;
  * `apt_invalid_operation` &ndash; fixes invalid `apt` and `apt-get` calls, like `apt-get isntall vim`;
  * `apt_list_upgradable` &ndash; helps you run `apt list --upgradable` after `apt update`;
  * `apt_upgrade` &ndash; helps you run `apt upgrade` after `apt list --upgradable`;
  * `aws_cli` &ndash; fixes misspelled commands like `aws dynamdb scan`;
  * `az_cli` &ndash; fixes misspelled commands like `az providers`;
  * `brew_cask_dependency` &ndash; installs cask dependencies;
  * `brew_install` &ndash; fixes formula name for `brew install`;
  * `brew_link` &ndash; adds `--overwrite --dry-run` if linking fails;
  * `brew_reinstall` &ndash; turns `brew install <formula>` into `brew reinstall <formula>`;
  * `brew_uninstall` &ndash; adds `--force` to `brew uninstall` if multiple versions were installed;
  * `brew_unknown_command` &ndash; fixes wrong brew commands, for example `brew docto/brew doctor`;
  * `brew_update_formula` &ndash; turns `brew update <formula>` into `brew upgrade <formula>`;
  * `cargo` &ndash; runs `cargo build` instead of `cargo`;
  * `cargo_no_command` &ndash; fixes wrong commands like `cargo buid`;
  * `cat_dir` &ndash; replaces `cat` with `ls` when you try to `cat` a directory;
  * `cd_correction` &ndash; spellchecks and corrects failed `cd` commands;
  * `cd_cs` &ndash; changes `cs` to `cd`;
  * `cd_mkdir` &ndash; creates directories before `cd`'ing into them;
  * `cd_parent` &ndash; changes `cd..` to `cd ..`;
  * `chmod_x` &ndash; adds execution bit;
  * `choco_install` &ndash; appends common suffixes for chocolatey packages;
  * `composer_not_command` &ndash; fixes composer command name;
  * `conda_mistype` &ndash; fixes conda commands;
  * `cp_create_destination` &ndash; creates a new directory when you attempt to `cp` or `mv` to a non-existent one;
  * `cp_omitting_directory` &ndash; adds `-a` when you `cp` a directory;
  * `cpp11` &ndash; adds missing `-std=c++11` to `g++` or `clang++`;
  * `dirty_untar` &ndash; fixes `tar x` command that untarred in the current directory;
  * `dirty_unzip` &ndash; fixes `unzip` command that unzipped in the current directory;
  * `django_south_ghost` &ndash; adds `--delete-ghost-migrations` to failed django south migration;
  * `django_south_merge` &ndash; adds `--merge` to inconsistent django south migration;
  * `docker_image_being_used_by_container` &ndash; removes the container that is using the image before removing the image;
  * `docker_login` &ndash; executes a `docker login` and repeats the previous command;
  * `docker_not_command` &ndash; fixes wrong docker commands like `docker tags`;
  * `dry` &ndash; fixes repetitions like `git git push`;
  * `fab_command_not_found` &ndash; fixes misspelled fabric commands;
  * `fix_alt_space` &ndash; replaces Alt+Space with Space character;
  * `fix_file` &ndash; opens a file with an error in your `$EDITOR`;
  * `gem_unknown_command` &ndash; fixes wrong `gem` commands;
  * `git_add` &ndash; fixes *"pathspec 'foo' did not match any file(s) known to git."*;
  * `git_add_force` &ndash; adds `--force` to `git add <pathspec>...` when paths are .gitignore'd;
  * `git_bisect_usage` &ndash; fixes `git bisect strt`, `git bisect goood`, etc. when bisecting;
  * `git_branch_0flag` &ndash; fixes commands such as `git branch 0v` and `git branch 0r`;
  * `git_branch_delete` &ndash; changes `git branch -d` to `git branch -D`;
  * `git_branch_delete_checked_out` &ndash; changes `git branch -d` to `git checkout master && git branch -D` when deleting a checked out branch;
  * `git_branch_exists` &ndash; offers `git branch -d foo`, `git branch -D foo` or `git checkout foo` when creating a branch that already exists;
  * `git_branch_list` &ndash; catches `git branch list` in place of `git branch` and removes created branch;
  * `git_checkout` &ndash; fixes branch name or creates new branch;
  * `git_clone_git_clone` &ndash; replaces `git clone git clone ...` with `git clone ...`;
  * `git_clone_missing` &ndash; adds `git clone` to URLs that appear to link to a git repository;
  * `git_commit_add` &ndash; offers `git commit -a ...` or `git commit -p ...` after failed commit with nothing staged;
  * `git_commit_amend` &ndash; offers `git commit --amend` after previous commit;
  * `git_commit_reset` &ndash; offers `git reset HEAD~` after previous commit;
  * `git_diff_no_index` &ndash; adds `--no-index` to `git diff` on untracked files;
  * `git_diff_staged` &ndash; adds `--staged` to `git diff` with unexpected output;
  * `git_fix_stash` &ndash; fixes `git stash` commands (misspelled subcommand and missing `save`);
  * `git_flag_after_filename` &ndash; fixes `fatal: bad flag '...' after filename`;
  * `git_go` &ndash; fixes `git install`/`git get`/`git mod` &rarr; `go install`/`go get`/`go mod`;
  * `git_help_aliased` &ndash; fixes `git help <alias>` replacing the alias with the aliased command;
  * `git_hook_bypass` &ndash; adds `--no-verify` flag to `git am`, `git commit`, or `git push`;
  * `git_lfs_mistype` &ndash; fixes mistyped `git lfs <command>` commands;
  * `git_main_master` &ndash; fixes incorrect branch name between `main` and `master`;
  * `git_merge` &ndash; adds remote to branch names;
  * `git_merge_unrelated` &ndash; adds `--allow-unrelated-histories` when required;
  * `git_not_command` &ndash; fixes wrong git commands like `git brnch`;
  * `git_pull` &ndash; sets upstream before executing `git pull`;
  * `git_pull_clone` &ndash; clones instead of pulling when the repo does not exist;
  * `git_pull_uncommitted_changes` &ndash; stashes changes before pulling and pops them afterwards;
  * `git_push` &ndash; adds `--set-upstream origin $branch` to failed `git push`;
  * `git_push_different_branch_names` &ndash; fixes pushes when local and remote branch names differ;
  * `git_push_force` &ndash; adds `--force-with-lease` to a `git push`;
  * `git_push_pull` &ndash; runs `git pull` when `push` was rejected;
  * `git_push_without_commits` &ndash; creates an initial commit when you `git add .` but forget to commit;
  * `git_rebase_merge_dir` &ndash; offers `git rebase (--continue | --abort | --skip)` or removing `.git/rebase-merge` when a rebase is in progress;
  * `git_rebase_no_changes` &ndash; runs `git rebase --skip` instead of `git rebase --continue` when there are no changes;
  * `git_remote_delete` &ndash; replaces `git remote delete` with `git remote remove`;
  * `git_remote_seturl_add` &ndash; runs `git remote add` when `git remote set-url` is used on a nonexistent remote;
  * `git_rm_local_modifications` &ndash; adds `-f` or `--cached` when you try to `rm` a locally modified file;
  * `git_rm_recursive` &ndash; adds `-r` when you try to `rm` a directory;
  * `git_rm_staged` &ndash; adds `-f` or `--cached` when you try to `rm` a file with staged changes;
  * `git_stash` &ndash; stashes your local modifications before rebasing or switching branch;
  * `git_stash_pop` &ndash; adds your local modifications before popping stash, then resets;
  * `git_tag_force` &ndash; adds `--force` to `git tag <tagname>` when the tag already exists;
  * `git_two_dashes` &ndash; adds a missing dash to commands like `git commit -amend` or `git rebase -continue`;
  * `go_git` &ndash; fixes `go pull`/`go push`/`go co` &rarr; `git pull`/`git push`/`git co`;
  * `go_run` &ndash; appends `.go` extension when compiling/running Go programs;
  * `go_unknown_command` &ndash; fixes wrong `go` commands, for example `go bulid`;
  * `gradle_no_task` &ndash; fixes not found or ambiguous `gradle` task;
  * `gradle_wrapper` &ndash; replaces `gradle` with `./gradlew`;
  * `grep_arguments_order` &ndash; fixes `grep` arguments order for situations like `grep -lir . test`;
  * `grep_recursive` &ndash; adds `-r` when you try to `grep` a directory;
  * `grunt_task_not_found` &ndash; fixes misspelled `grunt` commands;
  * `gulp_not_task` &ndash; fixes misspelled `gulp` tasks;
  * `has_exists_script` &ndash; prepends `./` when script/binary exists in the current directory;
  * `heroku_multiple_apps` &ndash; adds `--app <app>` to `heroku` commands like `heroku pg`;
  * `heroku_not_command` &ndash; fixes wrong `heroku` commands like `heroku log`;
  * `history` &ndash; replaces command with the most similar command from history;
  * `hostscli` &ndash; fixes `hostscli` usage;
  * `ifconfig_device_not_found` &ndash; fixes wrong device names like `wlan0` to `wlp2s0`;
  * `java` &ndash; removes `.java` extension when running Java programs;
  * `javac` &ndash; appends missing `.java` when compiling Java files;
  * `lein_not_task` &ndash; fixes wrong `lein` tasks like `lein rpl`;
  * `ln_no_hard_link` &ndash; catches hard link creation on directories, suggests symbolic link;
  * `ln_s_order` &ndash; fixes `ln -s` arguments order;
  * `long_form_help` &ndash; changes `-h` to `--help` when the short form is not supported;
  * `ls_all` &ndash; adds `-A` to `ls` when output is empty;
  * `ls_lah` &ndash; adds `-lah` to `ls`;
  * `man` &ndash; changes manual section;
  * `man_no_space` &ndash; fixes man commands without spaces, for example `mandiff`;
  * `mercurial` &ndash; fixes wrong `hg` commands;
  * `missing_space_before_subcommand` &ndash; fixes commands with missing space like `npminstall`;
  * `mkdir_p` &ndash; adds `-p` when you try to create a directory without a parent;
  * `mvn_no_command` &ndash; adds `clean package` to `mvn`;
  * `mvn_unknown_lifecycle_phase` &ndash; fixes misspelled lifecycle phases with `mvn`;
  * `nixos_cmd_not_found` &ndash; installs apps on NixOS;
  * `no_command` &ndash; fixes wrong console commands, for example `vom/vim`;
  * `no_such_file` &ndash; creates missing directories with `mv` and `cp` commands;
  * `npm_missing_script` &ndash; fixes `npm` custom script name in `npm run-script <script>`;
  * `npm_run_script` &ndash; adds missing `run-script` for custom `npm` scripts;
  * `npm_wrong_command` &ndash; fixes wrong npm commands like `npm urgrade`;
  * `omnienv_no_such_command` &ndash; fixes wrong commands for `goenv`, `nodenv`, `pyenv` and `rbenv`;
  * `open` &ndash; prepends `http://` to addresses passed to `open`, or creates missing files/directories;
  * `pacman` &ndash; installs app with `pacman` if it is not installed (uses `yay`, `pikaur` or `yaourt` if available);
  * `pacman_invalid_option` &ndash; replaces lowercase `pacman` options with uppercase;
  * `pacman_not_found` &ndash; fixes package name with `pacman`, `yay`, `pikaur` or `yaourt`;
  * `path_from_history` &ndash; replaces not-found path with a similar path from history;
  * `php_s` &ndash; replaces `-s` with `-S` when trying to run a local php server;
  * `pip_install` &ndash; fixes permission issues with `pip install` by adding `--user` or `sudo`;
  * `pip_unknown_command` &ndash; fixes wrong `pip` commands, for example `pip instatl`;
  * `port_already_in_use` &ndash; kills the process bound to a port;
  * `prove_recursively` &ndash; adds `-r` when called with a directory;
  * `python_command` &ndash; prepends `python` when you try to run a non-executable Python script;
  * `python_execute` &ndash; appends missing `.py` when executing Python files;
  * `python_module_error` &ndash; fixes `ModuleNotFoundError` by trying to `pip install` the module;
  * `quotation_marks` &ndash; fixes uneven usage of `'` and `"` in arguments;
  * `rails_migrations_pending` &ndash; runs pending migrations;
  * `react_native_command_unrecognized` &ndash; fixes unrecognized `react-native` commands;
  * `remove_shell_prompt_literal` &ndash; removes leading `$` prompt symbol, common when copying commands from documentation;
  * `remove_trailing_cedilla` &ndash; removes trailing cedillas `ç`, a common typo on European keyboards;
  * `rm_dir` &ndash; adds `-rf` when you try to remove a directory;
  * `rm_root` &ndash; adds `--no-preserve-root` to `rm -rf /`;
  * `scm_correction` &ndash; corrects wrong scm like `hg log` to `git log`;
  * `sed_unterminated_s` &ndash; adds missing `/` to `sed`'s `s` commands;
  * `sl` &ndash; changes `sl` to `ls`;
  * `ssh_known_hosts` &ndash; removes host from `known_hosts` on warning;
  * `sudo` &ndash; prepends `sudo` to the previous command if it failed due to permissions;
  * `sudo_command_from_user_path` &ndash; runs commands from user `$PATH` with `sudo`;
  * `switch_lang` &ndash; switches command from your local keyboard layout to en;
  * `systemctl` &ndash; correctly orders parameters of confusing `systemctl`;
  * `terraform_init` &ndash; runs `terraform init` before plan or apply;
  * `terraform_no_command` &ndash; fixes unrecognized `terraform` commands;
  * `tmux` &ndash; fixes `tmux` commands;
  * `touch` &ndash; creates missing directories before touching a file;
  * `tsuru_login` &ndash; runs `tsuru login` if not authenticated or session expired;
  * `tsuru_not_command` &ndash; fixes wrong `tsuru` commands like `tsuru shell`;
  * `unknown_command` &ndash; fixes hadoop hdfs-style "unknown command", e.g. adds missing `-` to `hdfs dfs ls`;
  * `unsudo` &ndash; removes `sudo` from previous command if the process refuses to run as superuser;
  * `vagrant_up` &ndash; starts up the vagrant instance;
  * `whois` &ndash; fixes `whois` command;
  * `workon_doesnt_exists` &ndash; fixes `virtualenvwrapper` env name or suggests creating a new one;
  * `wrong_hyphen_before_subcommand` &ndash; removes an improperly placed hyphen (`apt-install` &rarr; `apt install`);
  * `yarn_alias` &ndash; fixes aliased `yarn` commands like `yarn ls`;
  * `yarn_command_not_found` &ndash; fixes misspelled `yarn` commands;
  * `yarn_command_replaced` &ndash; fixes replaced `yarn` commands;
  * `yarn_help` &ndash; makes it easier to open `yarn` documentation;
  * `yum_invalid_operation` &ndash; fixes invalid `yum` calls, like `yum isntall vim`;

##### [Back to Contents](#contents)

## Contributing rules

Rules are static Go code in `internal/rules/`. Each rule is a file that calls
`register()` in its `init()` function:

```go
package rules

import "github.com/lyda/thefuck/internal/types"

func init() {
    register(Rule{
        Name: "my_rule",
        Match: func(cmd types.Command) bool {
            return strings.Contains(cmd.Output, "some error")
        },
        GetNewCommand: func(cmd types.Command) []types.CorrectedCommand {
            return single(strings.Replace(cmd.Script, "wrong", "right", 1))
        },
    })
}
```

Helper functions available in the `rules` package:

* `single(script)` &ndash; returns a single correction;
* `multi([]string)` &ndash; returns multiple corrections;
* `shellAnd(cmds...)` &ndash; joins commands with the shell's AND operator;
* `replaceArgument(script, from, to)` &ndash; replaces a shell word in a command;
* `getCloseMatches(word, possibilities, cutoff)` &ndash; fuzzy match, cutoff 0.0–1.0;
* `getAllMatchedCommands(output, separators)` &ndash; extracts suggestions after a separator line.

`types.Command` has two fields: `Script` (the command string) and `Output`
(combined stdout+stderr from re-running it), plus `ScriptParts()` which splits
the script into shell words.

##### [Back to Contents](#contents)

## Developing

```bash
git clone https://github.com/lyda/thefuck
cd thefuck
make check   # fmt, vet, staticcheck, gosec, test
make build   # builds ./bin/thefuck
make         # does both
```

## License MIT
Project License can be found [here](LICENSE.md).

[workflow-badge]:  https://github.com/lyda/thefuck/workflows/Tests/badge.svg
[workflow-link]:   https://github.com/lyda/thefuck/actions?query=workflow%3ATests
[license-badge]:   https://img.shields.io/badge/license-MIT-007EC7.svg
[examples-link]:   https://raw.githubusercontent.com/nvbn/thefuck/master/example.gif

##### [Back to Contents](#contents)
