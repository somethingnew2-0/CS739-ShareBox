from django.conf.urls import patterns, include, url
from django.contrib import admin
from server.views import *

urlpatterns = patterns('',
    # Examples:
    # url(r'^$', 'sharebox.views.home', name='home'),
    # url(r'^blog/', include('blog.urls')),

    url(r'^admin/', include(admin.site.urls)),
    ## User actions
    url(r'^user/new$', createUser),

    ## Client actions
    url(r'^client/(.*)/status$', getClientInitStatus),
    url(r'^client/(.*)/init$', initClient),
    url(r'^client/(.*)/file/add$', addFile),
)
