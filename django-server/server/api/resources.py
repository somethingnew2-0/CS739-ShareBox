## Not needed for now - we are not using Tastypie at the moment
## TODO : File should be removed later

from tastypie import fields
from tastypie.authentication import Authentication
from tastypie.authorization import Authorization
from tastypie.bundle import Bundle
from tastypie.exceptions import NotFound
from tastypie.resources import Resource
from models import *

#Non ORM data usage - adapted from,
#https://gist.github.com/nomadjourney/794424

class UserResource(Resource):
    # fields must map to the attributes in the User class
    id = fields.CharField(attribute = 'id')
    name = fields.CharField(attribute = 'name')
    files = fields.CharField(attribute = 'files')
    
    class Meta:
        resource_name = 'user'
        object_class = User
        authentication = Authentication()
        authorization = Authorization()
        consul_session = consulate.Consulate()
        consul_key_root = self._meta.resource_name + '/'
 
    # adapted this from ModelResource
    def get_resource_uri(self, bundle_or_obj):
        kwargs = {
            'resource_name': self._meta.resource_name,
        }
 
        if isinstance(bundle_or_obj, Bundle):
            kwargs['pk'] = bundle_or_obj.obj.id # pk is referenced in ModelResource
        else:
            kwargs['pk'] = bundle_or_obj.id
        
        if self._meta.api_name is not None:
            kwargs['api_name'] = self._meta.api_name
        
        return self._build_reverse_url('api_dispatch_detail', kwargs = kwargs)
 
    def get_object_list(self, request):
        # inner get of object list... this is where you'll need to
        # fetch the data from what ever data source
        return consul_session.kv.find(self._meta.consul_key_root)
 
    def obj_get_list(self, request = None, **kwargs):
        # outer get of object list... this calls get_object_list and
        # could be a point at which additional filtering may be applied
        return self.get_object_list(request)
 
    def obj_get(self, request = None, **kwargs):
        # get one object from data source
        pk = int(kwargs['pk'])
        try:
            return consul_session.kv[self._meta.consul_key_root + pk]
        except AttributeError:
            raise NotFound("Object not found") 
    
    def obj_create(self, bundle, request = None, **kwargs):
        # create a new row
        bundle.obj = User()
        
        # full_hydrate does the heavy lifting mapping the
        # POST-ed payload key/values to object attribute/values
        bundle = self.full_hydrate(bundle)
        
        # we add it to consult
        consul_session.kv[self._meta.consul_key_root + bundle.obj.id] = bundle.obj.__dict__
        return bundle
    
    def obj_update(self, bundle, request = None, **kwargs):
        # update an existing row
        pk = int(kwargs['pk'])
        try:
            bundle.obj = consul_session.kv[self._meta.consul_key_root + pk]
        except AttributeError:
            raise NotFound("Object not found")
        
        # let full_hydrate do its work
        bundle = self.full_hydrate(bundle)
        
        # update existing row in data dict
        consul_session.kv[self._meta.consul_key_root + pk] = bundle.obj
        return bundle