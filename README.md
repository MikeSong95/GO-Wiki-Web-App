Wiki web applications using GO.

Following this tutorial: https://golang.org/doc/articles/wiki/

Features added after tutorial completion:

1. Store templates in tmpl/ and page data in data/.
2. Add a handler to make the web root redirect to /view/FrontPage.
3. Spruce up the page templates by making them valid HTML and adding some CSS rules.
4. Implement inter-page linking by converting instances of [PageName] to 
<a href="/view/PageName">PageName</a>. 