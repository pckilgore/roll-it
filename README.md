# roll-it

Keep your AWS Credentials fresh 🍊 on Windows, Mac, Linux (arm or x86)!

Programmatically roll your AWS IAM access keys with rollit.

Detects your current credential chain and attempts to match it to a credentials file.

If successful, creates remote keys and deletes the old ones, updating your credentials
file.

This program fails safe: It will not delete any keys until is has successfully created 
new ones. Thus, you may need to delete old keys in the console before first using it
so that there is a free slot for your new key to be created in.

## Running Automatically

You need to use some other program to schedule execution.

For example, in a user crontab on Apple/Linux, extract the binary to somewhere in your
path, and add it to your user crontab:

```console
$ crontab -e # Edit to add...
$ crontab -l # Show what was added
50 12 * * 1-5 ~/bin/rollit
```
That example runs `rollit` Mon-Fri at 12:55 (assuming $HOME/bin is in your path).  

For more ideas, please [use the wiki](https://github.com/pckilgore/roll-it/wiki)

(c) 2021 Patrick Kilgore

See LICENSE
