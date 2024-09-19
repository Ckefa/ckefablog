package models

// Package IDs
var (
	DefaultPack  = 1
	BasicPack    = 2
	StandardPack = 3
	PremiumPack  = 4
)

// OrderDetails contains the features for each package
var OrderDetails = map[int][]string{
	DefaultPack: {
		"NO package Selected!!!",
	},
	BasicPack: {
		"✔ Responsive Design",
		"✔ 3 Custom Pages",
		"✔ Basic SEO Optimization",
		"✔ Social Media Integration",
		"✔ Custom Contact Form",
		"✔ 1 Month Support",
	},
	StandardPack: {
		"✔ 5 Custom Pages",
		"✔ Responsive Design",
		"✔ SEO Optimization",
		"✔ Database Integration",
		"✔ Content Management System (CMS)",
		"✔ Mobile Optimization",
		"✔ Image Gallery",
		"✔ 3 Months Support",
	},
	PremiumPack: {
		"✔ 10 Custom Pages",
		"✔ Responsive Design",
		"✔ Advanced SEO Optimization",
		"✔ Advanced Security Features (SSL, Firewalls)",
		"✔ Custom API Integration",
		"✔ Speed Optimization",
		"✔ E-Commerce Integration",
		"✔ Website Hosting and Maintenance",
		"✔ 6 Months Support",
		"✔ UX/UI Design Consultation",
		"✔ Content Strategy Consultation",
		"✔ Dedicated Server Hosting",
	},
}

// Additional functions or logic to interact with OrderDetails can be added here.
