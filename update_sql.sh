mv puamweb-1-12.sql prev_puamweb-1-12.sql
curl https://s3-us-west-2.amazonaws.com/puampretest/puamweb-1-12.sql > puamweb-1-12.sql
mysql -u root -phelloworld puamapi < puamweb-1-12.sql
