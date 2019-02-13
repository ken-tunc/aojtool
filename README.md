# aojtool

A CLI automation tool for [Aizu Online Jude](http://judge.u-aizu.ac.jp/onlinejudge).

- Run your code with sample inputs and compare it with the correct outputs
- Submit your code to AOJ
- Enable you to check your profile and recent submission status

## Installation

```console
$ go get github.com/ken-tunc/aojtool
```

Make sure that `$GPATH/bin` in your `$PATH`.

## Usage

### Login

Login to Aizu Online Judge(AOJ).
```console
$ aojtool login
```

### Run Code

Run your code with sample inputs and compare it with the correct outputs.
```console
$ aojtool run ALDS1_1_A main.c
```

You can specify language with `-l` flag,
defaults to your default programming language of AOJ.

### Submit Your Code

Submit your code to AOJ.
```console
$ aojtool submit ALDS1_1_A main.c
```

You can specify language with `-l` flag,
defaults to your default programming language of AOJ.

### Check Status

Check your profile and recent submission status.
```console
$ aojtool status
```

You can specify the number of submissions with `-n` flag.

### Logout

Logout from AOJ.
```console
$ aojtool logout
```

## License

See [LICENSE](LICENSE).
