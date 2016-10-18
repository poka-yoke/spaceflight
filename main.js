var child_process = require('child_process');

exports.handler = function(event, context) {
    var proc = child_process.spawn('./elasticsearch_client', [ JSON.stringify(event) ]/*, { stdio: 'inherit' }*/);
    
    proc.stdout.on('data', function (data) {    // stdout handler
  	    console.log('stdout: ' + data);
        context.succeed('stdout: ' + data);
	  });

    proc.stderr.on('data', function (data) {	// stderr handler
	      console.log('stderr: ' + data);
        context.succeed('stderr: ' + data);
        // context.fail('Something went wrong');
    });

    proc.on('exit', function (code) {
        console.log('lambda exited with code ' + code);
    });
};

