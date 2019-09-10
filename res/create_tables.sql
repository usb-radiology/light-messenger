CREATE TABLE `ArduinoStatus` (
  `departmentId` varchar(255) NOT NULL,
  `statusAt` bigint NOT NULL,
  PRIMARY KEY (`departmentId`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

CREATE TABLE `Notification` (
  `notificationId` varchar(255) NOT NULL,
  `modality` varchar(255) NOT NULL,
  `departmentId` varchar(255) NOT NULL,
  `priority` int(11) NOT NULL,
  `createdAt` bigint NOT NULL,
  `confirmedAt` bigint NOT NULL DEFAULT -1,
  `cancelledAt` bigint NOT NULL DEFAULT -1,
  PRIMARY KEY (`notificationId`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;
