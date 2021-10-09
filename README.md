# roll-it

Programmatically roll your AWS IAM access keys.

Detects your current credential chain and attempts to match it to a credentials file.

If successful, creates remote keys and deletes the old ones.

This program fails safe: It will not delete any keys until is has successfully created 
new ones. Thus, you may need to delete old keys in the console before first using it
so that there is a free slot for your new key to be created in.

(c) 2021 Patrick Kilgore
See LICENSE
