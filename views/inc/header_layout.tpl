<!DOCTYPE html>

<html>
    <head>
      <title>独孤影 - {{{.title}}}</title>
      <meta name="viewport" content="width=device-width, initial-scale=1"/>
      <meta http-equiv="Content-Type" content="text/html; charset=utf-8">
      <meta content="独孤影,博客,{{{.keywords}}}" name="keywords" />
      <meta content="{{{.description}}}" name="description" />
      <link rel="EditURI" type="application/rsd+xml" title="RSD" href="{{{.host}}}/xmlrpc" />
      <link rel="shortcut icon" href="/favicon.ico" />
      {{{if .inDev}}}
          {{{template "inc/css_dev.tpl" .}}}
      {{{else}}}
          {{{template "inc/css_prod.tpl" .}}}
      {{{end}}}
      <meta name="google-site-verification" content="ohMjRPHv0sKAahvl1H0GC7Dx0-z-zXbMNnWBfxp2PYY" />
      <meta name="baidu-site-verification" content="h3Y69jNgBz" />
    </head>
    <body >
      <div class="main">

          