using System;
using System.Threading;
using budgetride.AutomationResource;
using budgetride.Pages;
using Microsoft.VisualStudio.TestTools.UnitTesting;
using OpenQA.Selenium;

namespace budgetride.Test
{
	[TestClass]
	public class TestBase
	{
		public IWebDriver Driver { get; private set; }


		[TestInitialize]
		public void SetupForEverySingleTestMethod()
		{
			var factory = new WebDriverFactory();
			Driver = factory.Create(BrowserType.Chrome);
		}

  

		[TestCleanup]
		public void CleanUpAfterEveryTestMethod()
		{
			Driver.Close();
			Driver.Quit();
		}


	}
}
