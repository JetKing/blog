			<div class="header">
        		<div class="icons">
        			<a href="http://my.oschina.net/duguying" target="_black"><span title="follow me on oschina" class="icon-osc imgicon"></span></a>
              <a href="https://github.com/duguying" target="_black"><span title="follow me on Github" class="icon-github imgicon"></span></a>
        			<a href="http://weibo.com/duguying2008" target="_black"><span title="find me on Weibo" class="icon-weibo imgicon"></span></a>
        			<a href="http://gplus.to/duguying" target="_black"><span title="find me on g+" class="icon-gplus imgicon"></span></a>
        			<a href="https://twitter.com/duguying" target="_black"><span title="find me on Twitter" class="icon-twitter imgicon"></span></a>
        		</div>
        		<ul class="menu">
        			<li id="about">
                <a>关于</a>
                <div class="drop-menu">
                  <span class="droplist-array-down">◆</span>
                  <ul class="drop-list">
                    <li><a target="_blank" href="">关于博客</a></li>
                    <li><a target="_blank" href="/about/resume">个人简历</a></li>
                    <li><a target="_blank" href="/about/statistics">代码统计</a></li>
                  </ul>
                </div>
              </li>
        			<li>
                <a href="/project">项目</a>
              </li>
        			<li><a href="/list">列表</a></li>
        			<li><a href="/">博文</a></li>
        		</ul>
        		<div class="banner">
        			<a href="/" title="独孤影"><span class="title"><img src="/static/theme/default/img/dgy.svg" alt="独孤影"></span></a>
        		</div>
        		<div class="gap">
              {{{if .userIs}}}
              {{{if eq .userIs "admin"}}}
        			<a href="/admin" title="管理页面">
                <img class="gravatar" src="/logo" />
              </a>
              {{{else}}}
              <img class="gravatar" src="/logo" />
              {{{end}}}
              {{{end}}}
        		</div>
        	</div>
