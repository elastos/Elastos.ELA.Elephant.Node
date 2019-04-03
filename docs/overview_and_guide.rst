System Requirements and Installation Guide
########################################

.. toctree::
  :maxdepth: 2

Installation
============
The preferred way to install elephant node is to get the source code and build by yourself.

Mac
***

Check OS version,Make sure the OSX version is 16.7+::

    $ uname -srm
    Darwin 16.7.0 x86_64

Install Homebrew::

    /usr/bin/ruby -e "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/master/install)"

Install Git::

    brew install git

Install Go distribution 1.10 , Use Homebrew to install Golang 1.10::

    $ brew install go@1.10

Install Glide,Glide is a package manager for Golang. We use Glide to install dependent packages::

    $ brew install --ignore-dependencies glide

Ubuntu
******

Check OS version,Make sure your ubuntu version is 16.04+::

    $ cat /etc/issue
    Ubuntu 16.04.3 LTS \n \l

Install git::

    $ sudo apt-get install -y git

Install Go distribution 1.10::

    $ sudo apt-get install -y software-properties-common
    $ sudo add-apt-repository -y ppa:gophers/archive
    $ sudo apt update
    $ sudo apt-get install -y golang-1.10-go

Install Glide, Glide is a package manager for Golang. We use Glide to install dependent packages::

    $ cd ~/dev
    $ curl https://glide.sh/get | sh


Configure the Elastos.ELA::

    {
      "Configuration": {
        "ActiveNet": "testnet",
        //Local web portal port number. User can go to http://127.0.0.1:10333/info to access the web UI
        "HttpInfoPort": 21333,
        //true to start the webUI, false to disable
        "HttpInfoStart": true,
        //Restful port number
        "HttpRestPort": 21334,
        //Websocket port number
        "HttpWsPort": 21335,
        //RPC port number
        "HttpJsonPort": 21336,
        //P2P port number
        "NodePort": 21338,
        //Log level. Level 0 is the highest, 6 is the lowest.
        "PrintLevel": 1,
        "PowConfiguration": {
          //Mining reward receiving address
          "PayToAddr": "EeEkSiRMZqg5rd9a2yPaWnvdPcikFtsrjE",
          "MinerInfo": "ELA",
          //Minimal mining fee
          "MinTxFee": 100
        },
        "RpcConfiguration": {
          //User name: if set, you need to provide user name and password when calling the rpc interface
          "User": "ElaUser",
          //User password: if set, you need to provide user name and password when calling the rpc interface
          "Pass": "Ela123",
          //If hanve "0.0.0.0" in WhiteIPList will allow all ip to connect, otherwise only allow ip in WhiteIPList to connect
          "WhiteIPList": [
            "127.0.0.1"
          ]
        }
      }
    }

Extra feature configure::

    {
      //Get CMC Apikey go to https://coinmarketcap.com/api/
      "Cmc":{
        "ApiKey":["5c8bcab7-a811-428d-9c2b-3f0326de4f66","237bd580-3c68-42f8-b9e8-fd201c4933ac"],
        "Inteval":"15m"
      }
    }

.. warn::
    At the moment, elephant node only fire a Elastos.ELA node.

Build the node
**************

Setup basic workspace,In this instruction we use ~/dev/src/github.com/elastos as our working directory. If you clone the source code to a different directory, please make sure you change other environment variables accordingly (not recommended)::

    $ mkdir -p ~/dev/bin
    $ mkdir -p ~/dev/src/github.com/elastos/

Set correct environment variables::

    export GOROOT=/usr/local/opt/go@1.10/libexec
    export GOPATH=$HOME/dev
    export GOBIN=$GOPATH/bin
    export PATH=$GOROOT/bin:$PATH
    export PATH=$GOBIN:$PATH

Check Go version and glide version,Check the golang and glider version. Make sure they are the following version number or above::

    $ go version
    go version go1.10.2 darwin/amd64

    $ glide --version
    glide version 0.13.1

If you cannot see the version number, there must be something wrong when install.

Clone source code to $GOPATH/src/github/elastos folder,Make sure you are in the folder of $GOPATH/src/github.com/elastos::

    $ git clone https://github.com/elastos/Elastos.ELA.Elephant.Node.git

If clone works successfully, you should see folder structure like $GOPATH/src/github.com/elastos/Elastos.ELA.Elephant.Node/Makefile

Install dependencies using Glide::

    $ cd $GOPATH/src/github.com/elastos/Elastos.ELA.Elephant.Node
    $ glide update && glide install

Make,Build the node::

    $ cd $GOPATH/src/github.com/elastos/Elastos.ELA.Elephant.Node
    $ make

If you did not see any error message, congratulations, you have made the Elephant full node.

Run the node::

    $ ./elephant
