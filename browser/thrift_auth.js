/**
 * Orverride getXmlHttpRequestObject to set Authentication Token Header
 */
Thrift.TXHRAuthTransport = function() {
  var url = arguments[0]
  var options = arguments[1];
  Thrift.TXHRTransport.call(this, url, options);
  this.authToken = options.token;
  console.log("TXHRAuthTransport", this);
}
Object.setPrototypeOf(Thrift.TXHRAuthTransport.prototype, Thrift.TXHRTransport.prototype);

Thrift.TXHRAuthTransport.prototype.getXmlHttpRequestObject = function() {
  console.log("getXmlHttpRequestObject")
  var req = Thrift.TXHRTransport.prototype.getXmlHttpRequestObject.call(this);
  return req;
}
