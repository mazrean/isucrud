# DB Graph
node: ![](https://via.placeholder.com/16/795548/FFFFFF/?text=%20) `テーブル` ![](https://via.placeholder.com/16/1976D2/FFFFFF/?text=%20) `関数` 

edge: ![](https://via.placeholder.com/16/CDDC39/FFFFFF/?text=%20) `INSERT` ![](https://via.placeholder.com/16/FF9800/FFFFFF/?text=%20) `UPDATE` ![](https://via.placeholder.com/16/F44336/FFFFFF/?text=%20) `DELETE` ![](https://via.placeholder.com/16/78909C/FFFFFF/?text=%20) `SELECT` ![](https://via.placeholder.com/16/BBDEFB/FFFFFF/?text=%20) `関数呼び出し` 
```mermaid
graph LR
  classDef table fill:#795548,fill-opacity:0.5
  classDef func fill:#1976D2,fill-opacity:0.5
  func:github.com/catatsuy/private-isu/webapp/golang.getIndex[getIndex]:::func --> func:github.com/catatsuy/private-isu/webapp/golang.getSessionUser[getSessionUser]:::func
  func:github.com/catatsuy/private-isu/webapp/golang.getIndex[getIndex]:::func --> func:github.com/catatsuy/private-isu/webapp/golang.makePosts[makePosts]:::func
  func:github.com/catatsuy/private-isu/webapp/golang.getIndex[getIndex]:::func --> table:posts[posts]:::table
  func:github.com/catatsuy/private-isu/webapp/golang.makePosts[makePosts]:::func --> table:comments[comments]:::table
  func:github.com/catatsuy/private-isu/webapp/golang.makePosts[makePosts]:::func --> table:users[users]:::table
  func:github.com/catatsuy/private-isu/webapp/golang.postIndex[postIndex]:::func --> func:github.com/catatsuy/private-isu/webapp/golang.getSessionUser[getSessionUser]:::func
  func:github.com/catatsuy/private-isu/webapp/golang.postIndex[postIndex]:::func --> table:posts[posts]:::table
  func:github.com/catatsuy/private-isu/webapp/golang.getLogin[getLogin]:::func --> func:github.com/catatsuy/private-isu/webapp/golang.getSessionUser[getSessionUser]:::func
  func:github.com/catatsuy/private-isu/webapp/golang.getRegister[getRegister]:::func --> func:github.com/catatsuy/private-isu/webapp/golang.getSessionUser[getSessionUser]:::func
  func:github.com/catatsuy/private-isu/webapp/golang.getPostsID[getPostsID]:::func --> func:github.com/catatsuy/private-isu/webapp/golang.getSessionUser[getSessionUser]:::func
  func:github.com/catatsuy/private-isu/webapp/golang.getPostsID[getPostsID]:::func --> func:github.com/catatsuy/private-isu/webapp/golang.makePosts[makePosts]:::func
  func:github.com/catatsuy/private-isu/webapp/golang.getPostsID[getPostsID]:::func --> table:posts[posts]:::table
  func:github.com/catatsuy/private-isu/webapp/golang.tryLogin[tryLogin]:::func --> table:users[users]:::table
  func:github.com/catatsuy/private-isu/webapp/golang.getSessionUser[getSessionUser]:::func --> table:users[users]:::table
  func:github.com/catatsuy/private-isu/webapp/golang.getPosts[getPosts]:::func --> func:github.com/catatsuy/private-isu/webapp/golang.makePosts[makePosts]:::func
  func:github.com/catatsuy/private-isu/webapp/golang.getPosts[getPosts]:::func --> table:posts[posts]:::table
  func:github.com/catatsuy/private-isu/webapp/golang.getAccountName[getAccountName]:::func --> func:github.com/catatsuy/private-isu/webapp/golang.getSessionUser[getSessionUser]:::func
  func:github.com/catatsuy/private-isu/webapp/golang.getAccountName[getAccountName]:::func --> func:github.com/catatsuy/private-isu/webapp/golang.makePosts[makePosts]:::func
  func:github.com/catatsuy/private-isu/webapp/golang.getAccountName[getAccountName]:::func --> table:comments[comments]:::table
  func:github.com/catatsuy/private-isu/webapp/golang.getAccountName[getAccountName]:::func --> table:posts[posts]:::table
  func:github.com/catatsuy/private-isu/webapp/golang.getAccountName[getAccountName]:::func --> table:users[users]:::table
  func:github.com/catatsuy/private-isu/webapp/golang.getAdminBanned[getAdminBanned]:::func --> func:github.com/catatsuy/private-isu/webapp/golang.getSessionUser[getSessionUser]:::func
  func:github.com/catatsuy/private-isu/webapp/golang.getAdminBanned[getAdminBanned]:::func --> table:users[users]:::table
  func:github.com/catatsuy/private-isu/webapp/golang.postAdminBanned[postAdminBanned]:::func --> func:github.com/catatsuy/private-isu/webapp/golang.getSessionUser[getSessionUser]:::func
  func:github.com/catatsuy/private-isu/webapp/golang.postAdminBanned[postAdminBanned]:::func --> table:users[users]:::table
  func:github.com/catatsuy/private-isu/webapp/golang.getImage[getImage]:::func --> table:posts[posts]:::table
  func:github.com/catatsuy/private-isu/webapp/golang.postRegister[postRegister]:::func --> func:github.com/catatsuy/private-isu/webapp/golang.getSessionUser[getSessionUser]:::func
  func:github.com/catatsuy/private-isu/webapp/golang.postRegister[postRegister]:::func --> table:users[users]:::table
  func:github.com/catatsuy/private-isu/webapp/golang.postRegister[postRegister]:::func --> table:users[users]:::table
  func:github.com/catatsuy/private-isu/webapp/golang.postLogin[postLogin]:::func --> func:github.com/catatsuy/private-isu/webapp/golang.getSessionUser[getSessionUser]:::func
  func:github.com/catatsuy/private-isu/webapp/golang.postLogin[postLogin]:::func --> func:github.com/catatsuy/private-isu/webapp/golang.tryLogin[tryLogin]:::func
  func:github.com/catatsuy/private-isu/webapp/golang.postComment[postComment]:::func --> func:github.com/catatsuy/private-isu/webapp/golang.getSessionUser[getSessionUser]:::func
  func:github.com/catatsuy/private-isu/webapp/golang.postComment[postComment]:::func --> table:comments[comments]:::table
  linkStyle 0,1,5,7,8,9,10,14,16,17,21,23,26,29,30,31 stroke:#BBDEFB,stroke-width:2px
  linkStyle 2,3,4,11,12,13,15,18,19,20,22,25,28 stroke:#78909C,stroke-width:2px
  linkStyle 6,27,32 stroke:#CDDC39,stroke-width:2px
  linkStyle 24 stroke:#FF9800,stroke-width:2px
```