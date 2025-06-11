
# generated with CMTools 2.2.7 b3036dd 2025-06-02

#
# Makefile for running pandoc on all Markdown docs ending in .md
#
PROJECT = dataset

PANDOC = $(shell which pandoc)

MD_PAGES = $(shell ls -1 *.md)

HTML_PAGES = $(shell ls -1 *.md | sed -E 's/\.md/\.html/g')

build: $(HTML_PAGES) $(MD_PAGES)

$(HTML_PAGES): $(MD_PAGES) .FORCE
	if [ -f $(PANDOC) ]; then $(PANDOC) --metadata title=$(basename $@) -s --to html5 -o $(basename $@).html \
		--lua-filter=../links-to-html.lua \
	    --template=page.tmpl $(basename $@).md ; fi

clean:
	@rm *.html

.FORCE:
