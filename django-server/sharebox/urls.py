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
    url(r'^client/(.*)/recover$', recoverClient),
    url(r'^client/(.*)/file/add$', addFile),
    url(r'^client/(.*)/file/remove$', removeFile),
    url(r'^client/(.*)/file/update$', updateFile),
    url(r'^file/(.*)/commit$', commitFile),
    url(r'^file/(.*)/delete$', deleteFile),
    url(r'^file/(.*)/download$', downloadFile),
    url(r'^shard/(.*)/validate$', validateShard),
    url(r'^shard/(.*)/invalidate$', invalidateShard),
)
