# Thesis template

## Overleaf

Use [Overleaf](https://www.overleaf.com/) to work on your diploma thesis doc.
Just upload the template and compile the project specifying a XeLaTeX compiler (2021)
and `thesis.tex` as the main document.

## Local Env using Docker

### Create `thesis` image

```
docker build -t thesis .
```

### Start a new container

```
docker run -it --rm --name thesis-box \
-h thesis-box \
-e UID=$(id -u) \
-e GID=$(id -g) \
-e DOCKER_GID=$(getent group docker | cut -d: -f3) \
-e HOME=$HOME \
-e USER=$USER \
-v /var/run/docker.sock:/var/run/docker.sock \
-v $(pwd):/thesis-template \
thesis \
bash --login
```

Alternatively, you can use the `Makefile` to start the container:

```
make enter-docker
```

The thesis template directory is mounted under `/thesis-template`:

```
cd /thesis-template
```

### Create diploma thesis PDF

```
make
```

### Remove produced artifacts (not .pdf)

```
make clean
```

### Remove all artifacts produced

```
make distclean
```

## VSCode with Docker

We suggest using the `LaTeX Workshop` extension of VSCode for local development
using VSCode.

Moreover you can setup VSCode to use the `thesis-box` container we provide for
building the project.

You can read more about this on:
https://intellexlab.medium.com/docker-based-latex-environment-for-quick-easy-write-up-using-vscode-795b921cc411


Procedure:
1. Install the `Remote - Containers` extension, see https://marketplace.visualstudio.com/items?itemName=ms-vscode-remote.remote-containers
   for instructions.
2. Start the `thesis-box` container, i.e., run `make enter-docker`.
3. Open VSCode.
4. In VSCode, head to left bottom corner, and click the relevant button to open
   the Remote Containers drop-down menu. TODO: Add screenshot.
5. Click `Attach to Running Container`, and select `thesis-box`. A new VSCode
   window will pop up.
6. Go to File > Open Folder and type `/thesis-template` to open the folder
   mounted in the container.
7. Go to extensions, search for `LaTeX Workshop` extension, and click `Install
   in Container thesis (thesis-box)`.
8. Ensure that the workspace settings located under `.vscode` dir are loaded
   from the VSCode, so that it uses the Makefile to build the project.



