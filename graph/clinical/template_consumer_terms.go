package clinical

const consumerTermsTemplate = `
<!DOCTYPE html>
<html lang="en">
  <head>
    <!-- Required meta tags -->
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1, shrink-to-fit=no">

    <!-- Bootstrap CSS -->
    <!-- Latest compiled and minified CSS -->
    <link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/bootstrap/3.4.1/css/bootstrap.min.css">

    <!-- jQuery library -->
    <script src="https://ajax.googleapis.com/ajax/libs/jquery/3.5.1/jquery.min.js"></script>

    <!-- Latest compiled JavaScript -->
    <script src="https://maxcdn.bootstrapcdn.com/bootstrap/3.4.1/js/bootstrap.min.js"></script>
    <link rel="stylesheet" href="base.css">
    <link rel="stylesheet" href="invalid.css">

    <title>Consumer Terms And Conditions</title>
  </head>
  <body>
    <div class="container">
      <div class="column justify-content-center">
      
        <h2 class="headline" style="text-align:center;">CONSUMER TERMS AND CONDITIONS</h2>
      
        {{ . }}
           
      </div>
    </div>    
  </body>
</html>
`
