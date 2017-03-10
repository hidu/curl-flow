<?php
for($i=0;$i<1000;$i++){
    $req=array();
    $req['url']="http://127.0.0.1/test.php";
    echo json_encode($req)."\n";
}
