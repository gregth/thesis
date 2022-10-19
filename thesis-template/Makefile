#######################################################
#
#  
# Vangelis Koukis <vkoukis@cslab.ece.ntua.gr>
# April 2008
# Updated October 2009, preparing final PhD thesis
#
#######################################################

LATEX ?= xelatex -interaction nonstopmode -file-line-error -synctex=1
BIBTEX ?= bibtex
SED ?= /bin/sed

all: thesis.pdf

clean:
	rm -f *log *aux *blg *bbl *dvi *toc thesis.out *ps *~

distclean: clean
	rm -f thesis.pdf

thesis.pdf: *.tex thesis.bbl
	$(LATEX) thesis.tex
	$(LATEX) thesis.tex

thesis.aux: thesis.tex
	$(LATEX) thesis.tex

thesis.bbl: thesis.aux
	$(BIBTEX) thesis.aux

thesis.pdfview: thesis.pdf
	"$(PDFVIEWER)" thesis.pdf

prepare-container:
	docker build -t thesis ../

enter-docker: prepare-container
	docker run -it --rm --name thesis-box \
	-h thesis-box \
	-e UID=$$(id -u) \
	-e GID=$$(id -g) \
	-e DOCKER_GID=$$(getent group docker | cut -d: -f3) \
	-e HOME=$$HOME \
	-e USER=$$USER \
	-v /var/run/docker.sock:/var/run/docker.sock \
	-v "$(PWD)/..":/diplom \
	thesis \
	bash --login
