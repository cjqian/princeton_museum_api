mv puamweb-1-12.sql prev_puamweb-1-12.sql
curl https://s3-us-west-2.amazonaws.com/puampretest/puamweb.sql > puamweb.sql
mysql -u root -phelloworld puamapi < puamweb.sql
