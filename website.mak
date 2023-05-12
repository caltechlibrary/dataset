#
# Makefile for running pandoc on all Markdown docs ending in .md
#
PROJECT = dataset

MD_PAGES = $(shell ls -1 *.md)

HTML_PAGES = $(shell ls -1 *.md | sed -E 's/\.md/.html/g')

MD_DOCS_PAGES = $(shell ls -1 docs/*.md)

HTML_DOCS_PAGES = $(shell ls -1 docs/*.md | sed -E 's/\.md/.html/g')

MD_HOWTO_PAGES = $(shell ls -1 how-to/*.md)

HTML_HOWTO_PAGES = $(shell ls -1 how-to/*.md | sed -E 's/\.md/.html/g')

MD_LIB_PAGES = $(shell ls -1 how-to/*.md)

HTML_LIB_PAGES = $(shell ls -1 how-to/*.md | sed -E 's/\.md/.html/g')


build: $(HTML_PAGES) $(MD_PAGES) $(HTML_DOCS_PAGES) $(MD_DOCS_PAGES) $(HTML_HOWTO_PAGES) $(MD_HOWTO_PAGES) pagefind
	@for FNAME in $(HTML_PAGES); do git add "$$FNAME"; done
	@for FNAME in $(HTML_DOCS_PAGES); do git add "$$FNAME"; done
	@for FNAME in $(HTML_HOWTO_PAGES); do git add "$$FNAME"; done
	@git commit -am 'website build process'

$(HTML_PAGES): $(MD_PAGES) .FORCE
	pandoc --metadata title=$(basename $@) -s --to html5 $(basename $@).md -o $(basename $@).html \
		--lua-filter=links-to-html.lua \
	    --template=page.tmpl
	@if [ $@ = "README.html" ]; then mv README.html index.html; git add index.html; fi

$(HTML_DOCS_PAGES): $(MD_DOCS_PAGES) .FORCE
	pandoc --metadata title=$(basename $@) -s --to html5 $(basename $@).md -o $(basename $@).html \
		--lua-filter=links-to-html.lua \
	    --template=docs.tmpl

$(HTML_HOWTO_PAGES): $(MD_HOWTO_PAGES) .FORCE
	pandoc --metadata title=$(basename $@) -s --to html5 "$(basename $@).md" -o "$(basename $@).html" \
		--lua-filter=links-to-html.lua \
	    --template=how-to.tmpl

$(HTML_LIB_PAGES): $(MD_LIB_PAGES) .FORCE
	pandoc --metadata title=$(basename $@) -s --to html5 $(basename $@).md -o $(basename $@).html \
		--lua-filter=links-to-html.lua \
	    --template=docs.tmpl
	@if [ $@ = "README.html" ]; then mv README.html index.html; git add index.html; fi


pagefind: .FORCE
	pagefind --verbose --exclude-selectors="nav,header,footer" --bundle-dir ./pagefind --source .
	git add pagefind

clean:
	@for FNAME in $(HTML_PAGES); do rm "$${FNAME}"; fi
	@for FNAME in $(HTML_DOCS_PAGES); do rm "$${FNAME}"; fi
	@for FNAME in $(HTML_HOWTO_PAGES); do rm "$${FNAME}"; fi


.FORCE:
