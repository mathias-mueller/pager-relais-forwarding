# Deployment using Ansible

## 0. Install ansible

https://docs.ansible.com/ansible/latest/installation_guide/intro_installation.html#pip-install

## 1. Add servers to deploy to your inventory

Add hosts to `inventory.yaml`

See: https://docs.ansible.com/ansible/latest/user_guide/intro_inventory.html

## 2. Check if servers are reachable

```shell
ansible all --list-hosts

ansible all -m ping

```

## 3. Build & Deploy service

```shell
ansible-playbook -i hosts main.yml
```