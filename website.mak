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


build: $(HTML_PAGES) $(MD_PAGES) $(HTML_DOCS_PAGES) $(MD_DOCS_PAGES) $(HTML_HOWTO_PAGES) $(MD_HOWTO_PAGES) pagefind

$(HTML_PAGES): $(MD_PAGES) .FORCE
	pandoc --metadata title=$(basename $@) -s --to html5 $(basename $@).md -o $(basename $@).html \
		--lua-filter=links-to-html.lua \
	    --template=page.tmpl
	@if [ $@ = "README.html" ]; then mv README.html index.html; fi

$(HTML_DOCS_PAGES): $(MD_DOCS_PAGES) .FORCE
	pandoc --metadata title=$(basename $@) -s --to html5 $(basename $@).md -o $(basename $@).html \
		--lua-filter=links-to-html.lua \
	    --template=docs.tmpl

$(HTML_HOWTO_PAGES): $(MD_HOWTO_PAGES) .FORCE
	pandoc --metadata title=$(basename $@) -s --to html5 "$(basename $@).md" -o "$(basename $@).html" \
		--lua-filter=links-to-html.lua \
	    --template=how-to.tmpl


pagefind: .FORCE
	pagefind --verbose --exclude-selectors="nav,header,footer" --bundle-dir ./pagefind --source .
	git add pagefind

clean:
	@if [ -f index.html ]; then rm *.html; fi

.FORCE:
