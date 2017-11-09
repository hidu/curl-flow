<?php
//build requests

for($i=0;$i<10000;$i++){
    $req=array();
    $req['url']='http://127.0.0.1:8088/test.php?i='.$i;
    
    $req['method']='post';
    $req['header']=array(
        'Content-Type'=>'application/x-www-form-urlencoded',
        'head-a'=>'head-v',
    );
    $req['playload']=http_build_query(array('id'=>$i,'now'=>time()));
    echo json_encode($req)."\n";
}
