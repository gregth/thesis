#!/usr/bin/env python3
# -*- coding: utf-8 -*-
#
# This file is part of Rok.
#
# Copyright © 2020 Arrikto Inc.  All Rights Reserved.

"""Entrypoint in rok-dev.
Create a new user in the Docker image linking the user information from
physical host to the container.
"""

from os import environ
import sys
import grp
import pwd
import subprocess

USERNAME = "user"
GROUP_NAME = "thesis-dev"
GROUP_DOCKER = "docker"
GROUP_SUDO = "sudo"
HOME = "/home/user"
BASH_PATH = "/bin/bash"
SUDO_PATH = "/etc/sudoers.d/"
DOCKER_SOCK = "/var/run/docker.sock"


def set_user(user, uid, group, home_dir):
    """Set user parameteres.
    Check if the user already exists otherwise
    create a new user, set the user home directory
    and the main/secondary group (group and GROUP_SUDO)
    """
    try:
        pwd.getpwnam(user)
    except KeyError:
        subprocess.run(["useradd", "--shell", BASH_PATH, "-u", str(uid),
                        "-o", "-c", "", "-m", "-g", group, "-d", home_dir,
                        user])
        subprocess.run(["usermod", "-aG", GROUP_SUDO, user])


def get_group_name(gid):
    """Get the group name.
    If GID exist, return the linked group name otherwise create
    a new group and return GROUP_NAME
    """
    try:
        group = grp.getgrgid(int(gid)).gr_name
    except KeyError:
        group = GROUP_NAME
        subprocess.run(["groupadd", GROUP_NAME, "-g", str(gid)])

    return group


def set_docker_group(user, docker_gid):
    """Set Docker group and assing the user to it.
    Check if GROUP_DOCKER exist, otherwise create
    a new group with GROUP_DOCKER and gid_docker
    """
    try:
        gr_mem = grp.getgrnam(GROUP_DOCKER).gr_mem
        if user not in gr_mem:
            subprocess.run(["usermod", "-aG", GROUP_DOCKER, user])
    except KeyError:
        subprocess.run(["groupadd", GROUP_DOCKER, "-g", str(docker_gid)])
        subprocess.run(["usermod", "-aG", GROUP_DOCKER, user])


def main():
    # Initialize the env variables if not shared by the user
    uid = environ.get("UID")
    gid = environ.get("GID")
    docker_gid = environ.get("DOCKER_GID")
    user = environ.get("USER", USERNAME)
    home_dir = environ.get("HOME", HOME)
    group = get_group_name(gid)

    # Check/set the new user
    set_user(user, uid, group, home_dir)

    # Creation docker group and assign it to the user
    set_docker_group(user, docker_gid)

    # Set the group GROUP_SUDO to be able to run sudo without authenticating
    sudoers = open(SUDO_PATH + user, "w")
    sudoers.write("%" + GROUP_SUDO + " ALL=(ALL) NOPASSWD:ALL")
    sudoers.close()

    # Run bash changing the current user from root to user
    subprocess.run(["gosu", user, BASH_PATH], cwd=home_dir)


if __name__ == "__main__":
    sys.exit(main())