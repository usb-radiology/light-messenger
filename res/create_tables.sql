CREATE TABLE `ArduinoStatus` (
  `departmentId` varchar(255) NOT NULL,
  `statusAt` bigint NOT NULL,
  PRIMARY KEY (`departmentId`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- 2016/10/27 15:35:25


CREATE TABLE `Notification` (
  `notificationId` varchar(255) NOT NULL,
  `departmentId` varchar(255) NOT NULL,
  `notificationStatus` varchar(255) NOT NULL,
  `priority` int(11) NOT NULL,
  `superiorDepartment` varchar(255) NOT NULL,
  `createdAt` varchar(255) NOT NULL,
  `processedAt` varchar(255)DEFAULT NULL,
  PRIMARY KEY (`notificationId`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;