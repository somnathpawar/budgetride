using System;
using System.Threading;
using OpenQA.Selenium;

namespace budgetride.Pages
{
    public class BaseApplicationPage
    {
		protected IWebDriver Driver { get; set; }
        public BaseApplicationPage(IWebDriver driver)
        {
            Driver = driver;
        }

        internal void GoTO()
        {
			Driver.Navigate().GoToUrl("https://tech9app.com/");
          
			Thread.Sleep(4000);
        }
    }
}
