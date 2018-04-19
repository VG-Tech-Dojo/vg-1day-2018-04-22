# -*- mode: ruby -*-
# vi: set ft=ruby :

Vagrant.configure("2") do |config|
  config.vm.box = "ubuntu/xenial64"
  config.vm.box_check_update = false

  config.vm.network "forwarded_port", guest: 8080, host: 8080, id: 'http'

  config.vm.synced_folder ".", "/home/vagrant/go/src/github.com/VG-Tech-Dojo/vg-1day-2018-04-22"

  config.vm.provision "shell", inline: <<-SHELL
    add-apt-repository ppa:gophers/archive
    curl -fsSL https://download.docker.com/linux/ubuntu/gpg | apt-key add -
    add-apt-repository "deb [arch=amd64] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable"
    apt-get update
    apt-get install -y git build-essential apt-transport-https ca-certificates curl software-properties-common
    apt-get install -y golang-1.10-go
    apt-get install -y docker-ce
    echo "export PATH=$PATH:/usr/lib/go-1.10/bin:/home/vagrant/go/bin" >> /home/vagrant/.bashrc
    echo "export GOPATH=/home/vagrant/go" >> /home/vagrant/.bashrc
    chown -R vagrant /home/vagrant/go
  SHELL
end
