# roll-it

Keep your AWS Credentials fresh 🍊 on Windows, Mac, Linux (arm or x86)!

## What it Does

Programmatically roll your AWS IAM access keys and secret. Helpful for
automatically ensuring compliance with organization key rotation policies.

Detects your current credential chain and attempts to match it to a credentials
file.

If successful, creates remote keys and deletes the old ones, updating your
credentials file.

This program fails safe: It will not delete any keys until is has successfully
created new ones. Thus, you may need to delete old keys in the console before
first using it so that there is a free slot for your new key to be created in.

## What it Does Not Do

Schedule itself. Use other tools for that.

## Installation

Download the latest [release](https://github.com/pckilgore/roll-it/releases),
for your OS and architecture, and either put it in your PATH or somewhere your
scheduler can access both it and your `~/.aws/credentials` file. There are no
options and no output for a successful key rotation, but a failure will print
information about the error and exit with a non-zero exit code.

## Project Goals

 1. Auditability:
    - No dependencies outside Go standard library and official AWS SDKs. Your
      associate engineers should be able to determine this program is safe to
      handle your organization's sensitive credentials.
 3. Respectful and Safe:
    - It should always fail-safe, and try to never leave the user without valid
      credentials, even if that means console intervention is neccessary for it
      to begin (or continue) to work.
 3. Sane defaults:
    - It should work out of the box for standard setups, and should not be
      extended to support exotic setups.
 4. Cross-Platform:
    - Make no assumptions about operating system or scheduler.
 5. Small & Fast:
    - You should never notice it running.
 6. User Facing:
    - Not designed—and should not be used—to rotate server or service keys.

## Running Automatically

You need to use some other program to schedule execution.

For example, in a user crontab on Apple/Linux, extract the binary to somewhere,
and add an entry to your user crontab referencing the executable:

```console
$ crontab -e # Edit to add...
$ crontab -l # Show what was added
50 9 * * 1-5 ~/bin/rollit
```
That example runs `rollit` Mon-Fri at 9:55 AM (assuming rollit executable copied
to $HOME/bin).

For more ideas, please [use the wiki](https://github.com/pckilgore/roll-it/wiki)

(c) 2021 Patrick Kilgore

See LICENSE
